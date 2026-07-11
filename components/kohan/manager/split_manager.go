package manager

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/rs/zerolog/log"
)

// SplitManager normalizes broker trades and carry-over accounts onto a
// synthetic split-adjusted share basis.  Quantities are scaled by the
// inferred split factor; unit prices are divided by the same factor.
// USDValue, commission and MarketValue remain economically invariant.
//
// The factor is inferred per (ticker, date) by comparing the broker
// unit price with the Yahoo adjusted close on the same date.  Only
// supported forward factors (2,3,4,5,6,8,10) and their reciprocals
// are auto-selected.  Mismatches that do not match a supported factor
// within 5% tolerance fail fast.
type SplitManager interface {
	// NormalizeTrades returns a copy of trades with adjusted quantity
	// and unit price.  USDValue and Commission are preserved.
	NormalizeTrades(ctx context.Context, trades []tax.Trade) ([]tax.Trade, common.HttpError)
	// NormalizeAccount returns a copy of the account with adjusted
	// quantity, MarketValue, OriginQty and OriginPrice.
	NormalizeAccount(ctx context.Context, account tax.Account) (tax.Account, common.HttpError)
}

// supportedFactors lists forward split ratios that can be auto-inferred.
// Reciprocals are implicitly supported for reverse splits.
var supportedFactors = []float64{2, 3, 4, 5, 6, 8, 10}

const factorTolerance = 0.05 // 5 %

type SplitManagerImpl struct {
	tickerManager TickerManager
}

func NewSplitManager(tickerManager TickerManager) *SplitManagerImpl {
	return &SplitManagerImpl{
		tickerManager: tickerManager,
	}
}

var _ SplitManager = (*SplitManagerImpl)(nil)

// ---------------------------------------------------------------------------
// NormalizeTrades
// ---------------------------------------------------------------------------

func (s *SplitManagerImpl) NormalizeTrades(ctx context.Context, trades []tax.Trade) ([]tax.Trade, common.HttpError) {
	result := make([]tax.Trade, len(trades))
	copy(result, trades)

	for i, t := range result {
		normalized, err := s.normalizeTrade(ctx, t)
		if err != nil {
			return nil, err
		}
		result[i] = normalized
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Factor inference
// ---------------------------------------------------------------------------

// inferFactor computes the split adjustment factor for a (brokerPrice, yahooClose) pair.
// Returns 1 when no split is detected.
func (s *SplitManagerImpl) inferFactor(brokerPrice, yahooPrice float64, _, _ string) float64 {
	if brokerPrice == 0 || yahooPrice == 0 {
		return 1
	}
	ratio := brokerPrice / yahooPrice

	// Ratio close to 1 → no split
	if s.matchesFactor(ratio, 1) {
		return 1
	}

	// Check supported forward factors
	for _, f := range supportedFactors {
		if s.matchesFactor(ratio, f) {
			return f
		}
	}

	// Check supported reverse factors (reciprocals)
	for _, f := range supportedFactors {
		if s.matchesFactor(ratio, 1/f) {
			return 1 / f
		}
	}

	// Ambiguous — return a non-factor sentinel
	return -1
}

// matchesFactor returns true when ratio is within tolerance of candidate.
func (s *SplitManagerImpl) matchesFactor(ratio, candidate float64) bool {
	if candidate == 0 {
		return false
	}
	diff := math.Abs(ratio-candidate) / candidate
	return diff <= factorTolerance
}

// ---------------------------------------------------------------------------
// findFactor wraps inferFactor and fails with an error when the mismatch
// approaches a known factor but falls outside tolerance, or when the
// observed ratio is large enough to suggest a corporate action but does
// not match any supported factor.  Ratios between 0.5 and 2.0 that are
// not close to a supported factor are treated as normal market variation
// (factor = 1) rather than causing a hard error.
func (s *SplitManagerImpl) findFactor(brokerPrice, yahooPrice float64, ticker, date string) (float64, common.HttpError) {
	factor := s.inferFactor(brokerPrice, yahooPrice, ticker, date)
	if factor >= 0 {
		return factor, nil
	}

	// Determine which candidate was closest so the error message is useful.
	ratio := brokerPrice / yahooPrice
	// Only fail on ratios outside [0.3, 3.0] — inside that range treat as
	// normal market variation even if no supported factor is an exact match.
	const minFailRatio = 0.3
	const maxFailRatio = 3.0
	if ratio < minFailRatio || ratio > maxFailRatio {
		return 0, common.NewHttpError(
			fmt.Sprintf("unsupported split ratio for %s on %s: broker=%.2f yahoo=%.2f ratio=%.3f",
				ticker, date, brokerPrice, yahooPrice, ratio),
			http.StatusBadRequest,
		)
	}
	// Normal market variation — treat as no split.
	return 1, nil
}

// make normalizeTrade use findFactor instead of inferFactor
func (s *SplitManagerImpl) normalizeTrade(ctx context.Context, trade tax.Trade) (tax.Trade, common.HttpError) {
	date, dateErr := time.Parse(time.DateOnly, trade.Date)
	if dateErr != nil {
		return trade, common.NewServerError(fmt.Errorf("invalid trade date %q: %w", trade.Date, dateErr))
	}

	yahooPrice, priceErr := s.tickerManager.GetPrice(ctx, trade.Symbol, date)
	if priceErr != nil {
		log.Warn().Str("Ticker", trade.Symbol).Str("Date", trade.Date).
			Err(priceErr).Msg("SplitManager: no Yahoo price for trade, left unchanged")
		return trade, nil
	}

	factor, err := s.findFactor(trade.USDPrice, yahooPrice, trade.Symbol, trade.Date)
	if err != nil {
		return trade, err
	}
	if factor == 1 {
		return trade, nil
	}

	log.Warn().
		Str("Ticker", trade.Symbol).
		Str("Date", trade.Date).
		Float64("observed_ratio", round2(trade.USDPrice/yahooPrice)).
		Float64("selected_factor", factor).
		Float64("original_qty", trade.Quantity).
		Float64("normalized_qty", round2(trade.Quantity*factor)).
		Float64("original_price", trade.USDPrice).
		Float64("normalized_price", round2(trade.USDPrice/factor)).
		Msg("SplitManager: inferred stock split — adjusting trade")

	return tax.Trade{
		Symbol:     trade.Symbol,
		Date:       trade.Date,
		Type:       trade.Type,
		Quantity:   round2(trade.Quantity * factor),
		USDPrice:   round2(trade.USDPrice / factor),
		USDValue:   trade.USDValue,
		Commission: trade.Commission,
	}, nil
}

func (s *SplitManagerImpl) NormalizeAccount(ctx context.Context, account tax.Account) (tax.Account, common.HttpError) {
	if account.OriginDate == "" {
		return account, nil
	}
	originDate, dateErr := time.Parse(time.DateOnly, account.OriginDate)
	if dateErr != nil {
		return account, common.NewServerError(fmt.Errorf("invalid origin date %q: %w", account.OriginDate, dateErr))
	}

	yahooPrice, priceErr := s.tickerManager.GetPrice(ctx, account.Symbol, originDate)
	if priceErr != nil {
		log.Warn().Str("Ticker", account.Symbol).Err(priceErr).Msg("SplitManager: no Yahoo price for account, left unchanged")
		return account, nil
	}

	factor, err := s.findFactor(account.OriginPrice, yahooPrice, account.Symbol, account.OriginDate)
	if err != nil {
		return account, err
	}
	if factor == 1 {
		return account, nil
	}

	log.Warn().
		Str("Ticker", account.Symbol).
		Str("OriginDate", account.OriginDate).
		Float64("observed_ratio", round2(account.OriginPrice/yahooPrice)).
		Float64("selected_factor", factor).
		Float64("original_qty", account.Quantity).
		Float64("normalized_qty", round2(account.Quantity*factor)).
		Msg("SplitManager: inferred stock split — adjusting account")

	return tax.Account{
		Symbol:      account.Symbol,
		Quantity:    round2(account.Quantity * factor),
		MarketValue: round2(account.MarketValue * factor),
		OriginDate:  account.OriginDate,
		OriginQty:   round2(account.OriginQty * factor),
		OriginPrice: round2(account.OriginPrice / factor),
	}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func round2(v float64) float64 {
	const decimalPlaces = 100.0
	return math.Round(v*decimalPlaces) / decimalPlaces
}
