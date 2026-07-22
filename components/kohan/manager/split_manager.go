package manager

import (
	"context"
	"fmt"
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

// tradeDateRange holds the inclusive date span for a single ticker's trades.
// Both values are pre-parsed time.Time values, ensuring deterministic
// downstream calls that depend only on trade data, not wall-clock time.
type tradeDateRange struct {
	earliest time.Time
	latest   time.Time
}

// tickerInfo groups trade indices, their pre-parsed dates, and the inclusive
// date range for a single ticker.  parsedDates[i] corresponds to indices[i].
type tickerInfo struct {
	indices     []int
	parsedDates []time.Time
	rng         tradeDateRange
}

// ---------------------------------------------------------------------------
// NormalizeTrades
// ---------------------------------------------------------------------------

func (s *SplitManagerImpl) NormalizeTrades(ctx context.Context, trades []tax.Trade) ([]tax.Trade, common.HttpError) {
	result := make([]tax.Trade, len(trades))
	copy(result, trades)
	if len(trades) == 0 {
		return result, nil
	}

	tickerMap, httpErr := validateAndGroupTrades(trades)
	if httpErr != nil {
		return nil, httpErr
	}

	for ticker, info := range tickerMap {
		if err := s.normalizeTickerTrades(ctx, ticker, info, result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// validateAndGroupTrades parses every trade date, validates it, and groups
// trades by ticker in a single pass.  Returns a map from ticker to tickerInfo
// that holds each ticker's global trade indices, parsed dates, and inclusive
// date range.  On the first invalid date a BadRequest error is returned and
// no TickerManager.GetSplits call is made.
func validateAndGroupTrades(trades []tax.Trade) (map[string]*tickerInfo, common.HttpError) {
	tickerMap := make(map[string]*tickerInfo)
	for i, t := range trades {
		d, err := time.Parse(time.DateOnly, t.Date)
		if err != nil {
			return nil, common.NewHttpError(
				fmt.Sprintf("invalid trade date %q for ticker %s", t.Date, t.Symbol),
				http.StatusBadRequest,
			)
		}
		info, ok := tickerMap[t.Symbol]
		if !ok {
			tickerMap[t.Symbol] = &tickerInfo{
				indices:     []int{i},
				parsedDates: []time.Time{d},
				rng: tradeDateRange{
					earliest: d,
					latest:   d,
				},
			}
		} else {
			if d.Before(info.rng.earliest) {
				info.rng.earliest = d
			}
			if d.After(info.rng.latest) {
				info.rng.latest = d
			}
			info.indices = append(info.indices, i)
			info.parsedDates = append(info.parsedDates, d)
		}
	}
	return tickerMap, nil
}

// normalizeTickerTrades fetches splits for a single ticker and adjusts all
// trades for that ticker in the result slice using pre-parsed dates held
// inside info.
func (s *SplitManagerImpl) normalizeTickerTrades(ctx context.Context, ticker string, info *tickerInfo, result []tax.Trade) common.HttpError {
	// Fetch splits once per ticker — from earliest to latest trade date
	splits, httpErr := s.tickerManager.GetSplits(ctx, ticker, info.rng.earliest, info.rng.latest)
	if httpErr != nil {
		return httpErr
	}

	if vErr := validateSplits(splits, ticker); vErr != nil {
		return vErr
	}

	// Normalize using pre-parsed dates from the grouped info
	for i, idx := range info.indices {
		result[idx] = s.applySplitsToTrade(result[idx], splits, info.parsedDates[i])
	}

	return nil
}

// ---------------------------------------------------------------------------
// Split application helpers
// ---------------------------------------------------------------------------

// applySplitsToTrade returns a copy of trade with quantity and price
// adjusted by the cumulative factor from all split events on strictly
// later UTC calendar dates.  Only Quantity and USDPrice are modified;
// USDValue, Commission, Symbol, Date and Type are preserved from the
// original.  refDay must be the pre-parsed calendar-day time.Time of
// the trade's date (caller validates the date before calling this).
func (s *SplitManagerImpl) applySplitsToTrade(trade tax.Trade, splits []tax.SplitInfo, refDay time.Time) tax.Trade {
	factor := cumulativeSplitFactor(splits, refDay)

	if factor == 1 {
		return trade
	}

	log.Warn().
		Str("Ticker", trade.Symbol).
		Str("Date", trade.Date).
		Float64("selected_factor", factor).
		Float64("original_qty", trade.Quantity).
		Float64("normalized_qty", trade.Quantity*factor).
		Float64("original_price", trade.USDPrice).
		Float64("normalized_price", trade.USDPrice/factor).
		Msg("SplitManager: event-based split adjustment — adjusting trade")

	result := trade
	result.Quantity = trade.Quantity * factor
	result.USDPrice = trade.USDPrice / factor
	return result
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// validateSplits validates every SplitInfo event in the slice for the given
// ticker. Returns the first validation error, or nil if all splits are valid.
func validateSplits(splits []tax.SplitInfo, ticker string) common.HttpError {
	for _, split := range splits {
		if vErr := split.Validate(ticker); vErr != nil {
			return vErr
		}
	}
	return nil
}

// cumulativeSplitFactor returns the product of all split ratios for events
// whose UTC calendar date is strictly later than refDay. The reference day
// is normalized to UTC midnight before comparison.
func cumulativeSplitFactor(splits []tax.SplitInfo, refDay time.Time) float64 {
	normalizedDay := refDay.UTC().Truncate(24 * time.Hour) //nolint:mnd
	factor := 1.0
	for _, split := range splits {
		if split.EffectiveDate().After(normalizedDay) {
			factor *= split.Ratio()
		}
	}
	return factor
}
