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
// synthetic split-adjusted share basis using explicit Yahoo Finance
// split events.  Quantities are scaled by the cumulative split factor;
// unit prices are divided by the same factor.  USDValue and Commission
// remain economically invariant.
//
// The factor is the product of (numerator/denominator) for every split
// event whose UTC calendar date is strictly later than the trade date.
// Split events are fetched once per ticker via TickerManager.GetSplits.
type SplitManager interface {
	// NormalizeTrades returns a copy of trades with adjusted quantity
	// and unit price.  USDValue and Commission are preserved.
	NormalizeTrades(ctx context.Context, trades []tax.Trade) ([]tax.Trade, common.HttpError)
}

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
	if len(trades) == 0 {
		return result, nil
	}

	// Group trade indices by ticker so we fetch splits once per ticker
	tickerIndices := make(map[string][]int)
	for i, t := range trades {
		tickerIndices[t.Symbol] = append(tickerIndices[t.Symbol], i)
	}

	for ticker, indices := range tickerIndices {
		if err := s.normalizeTickerTrades(ctx, ticker, indices, result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// normalizeTickerTrades fetches splits for a single ticker and adjusts all
// trades for that ticker in the result slice in-place.
func (s *SplitManagerImpl) normalizeTickerTrades(ctx context.Context, ticker string, indices []int, result []tax.Trade) common.HttpError {
	// Find earliest trade date for this ticker
	earliest := result[indices[0]].Date
	for _, idx := range indices[1:] {
		if result[idx].Date < earliest {
			earliest = result[idx].Date
		}
	}

	earliestDate, dateErr := time.Parse(time.DateOnly, earliest)
	if dateErr != nil {
		return common.NewServerError(
			fmt.Errorf("invalid trade date %q for ticker %s: %w", earliest, ticker, dateErr),
		)
	}

	// Fetch splits once per ticker — from earliest trade to current UTC date
	splits, httpErr := s.tickerManager.GetSplits(ctx, ticker, earliestDate, time.Now().UTC())
	if httpErr != nil {
		return httpErr
	}

	// Validate every split event once per ticker
	for _, split := range splits {
		if vErr := validateSplit(split, ticker); vErr != nil {
			return vErr
		}
	}

	// Parse and validate every trade date, then normalize
	for _, idx := range indices {
		tradeDate, dateErr := time.Parse(time.DateOnly, result[idx].Date)
		if dateErr != nil {
			return common.NewHttpError(
				fmt.Sprintf("invalid trade date %q for ticker %s", result[idx].Date, ticker),
				http.StatusBadRequest,
			)
		}
		tradeDay := calendarDay(tradeDate.Unix())
		result[idx] = s.applySplitsToTrade(result[idx], splits, tradeDay)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Split application helpers
// ---------------------------------------------------------------------------

// applySplitsToTrade returns a copy of trade with quantity and price
// adjusted by the cumulative factor from all split events on strictly
// later UTC calendar dates.  USDValue and Commission are unchanged.
// tradeDay must be the pre-parsed calendar-day Unix timestamp of the
// trade's date (caller validates the date before calling this).
func (s *SplitManagerImpl) applySplitsToTrade(trade tax.Trade, splits []tax.YahooSplit, tradeDay int64) tax.Trade {
	factor := 1.0
	for _, split := range splits {
		splitDay := calendarDay(split.Date)
		if splitDay > tradeDay {
			factor *= split.Numerator / split.Denominator
		}
	}

	if factor == 1 {
		return trade
	}

	log.Warn().
		Str("Ticker", trade.Symbol).
		Str("Date", trade.Date).
		Float64("selected_factor", factor).
		Float64("original_qty", trade.Quantity).
		Float64("normalized_qty", round2(trade.Quantity*factor)).
		Float64("original_price", trade.USDPrice).
		Float64("normalized_price", round2(trade.USDPrice/factor)).
		Msg("SplitManager: event-based split adjustment — adjusting trade")

	return tax.Trade{
		Symbol:     trade.Symbol,
		Date:       trade.Date,
		Type:       trade.Type,
		Quantity:   round2(trade.Quantity * factor),
		USDPrice:   round2(trade.USDPrice / factor),
		USDValue:   trade.USDValue,
		Commission: trade.Commission,
	}
}

// ---------------------------------------------------------------------------
// Validation
// ---------------------------------------------------------------------------

// validateSplit checks that a YahooSplit has valid event data:
// positive finite numerator and denominator, and a valid usable Unix timestamp.
// The error message includes the ticker name.
func validateSplit(split tax.YahooSplit, ticker string) common.HttpError {
	prefix := fmt.Sprintf("ticker %s: split event timestamp %d", ticker, split.Date)
	if split.Date <= 0 {
		return common.NewHttpError(fmt.Sprintf("%s: non-positive timestamp %d", prefix, split.Date), http.StatusBadRequest)
	}
	if split.Numerator <= 0 || math.IsInf(split.Numerator, 0) || math.IsNaN(split.Numerator) {
		return common.NewHttpError(fmt.Sprintf("%s: non-positive or non-finite numerator %f", prefix, split.Numerator), http.StatusBadRequest)
	}
	if split.Denominator <= 0 || math.IsInf(split.Denominator, 0) || math.IsNaN(split.Denominator) {
		return common.NewHttpError(fmt.Sprintf("%s: non-positive or non-finite denominator %f", prefix, split.Denominator), http.StatusBadRequest)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func round2(v float64) float64 {
	const decimalPlaces = 100.0
	return math.Round(v*decimalPlaces) / decimalPlaces
}
