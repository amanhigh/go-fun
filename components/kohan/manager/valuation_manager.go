package manager

import (
	"context"
	"time"

	"net/http"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type ValuationManager interface {
	// TODO: Build Broker Repository
	AnalyzeValuation(ctx context.Context, trades []tax.Trade, year int) (tax.Valuation, common.HttpError)
}

type ValuationManagerImpl struct {
	tickerManager TickerManager
}

func NewValuationManager(tickerManager TickerManager) ValuationManager {
	return &ValuationManagerImpl{
		tickerManager: tickerManager,
	}
}

func (v *ValuationManagerImpl) AnalyzeValuation(ctx context.Context, trades []tax.Trade, year int) (tax.Valuation, common.HttpError) {
	if len(trades) == 0 {
		return tax.Valuation{}, common.NewHttpError("no trades provided", http.StatusBadRequest)
	}

	// Validate all trades are for same ticker
	for _, t := range trades {
		if t.Symbol != trades[0].Symbol {
			return tax.Valuation{}, common.NewHttpError("multiple tickers found in trades", http.StatusBadRequest)
		}
	}

	// Initialize analysis
	analysis := tax.Valuation{
		Ticker: trades[0].Symbol,
	}

	// Track running position
	var currentPosition float64
	var maxPosition float64
	var firstBuyDate time.Time

	// Process trades chronologically
	for _, t := range trades {
		// Update position based on trade type
		if t.Type == "BUY" {
			currentPosition += t.Quantity
		} else {
			currentPosition -= t.Quantity
		}

		// Track first buy
		if firstBuyDate.IsZero() && t.Type == "BUY" {
			if parsedTime, err := t.GetDate(); err == nil {
				firstBuyDate = parsedTime
				analysis.FirstPosition = tax.Position{
					Date:     parsedTime,
					Quantity: t.Quantity,
					USDPrice: t.USDPrice,
				}
			}
		}

		// Track peak position
		// HACK: Handle Case of Peak TT Rate with changed Position
		if currentPosition > maxPosition {
			maxPosition = currentPosition
			if parsedTime, err := t.GetDate(); err == nil {
				analysis.PeakPosition = tax.Position{
					Date:     parsedTime,
					Quantity: currentPosition,
					USDPrice: t.USDPrice,
				}
			}
		}
	}

	// Get year end position if any holdings exist
	if currentPosition > 0 {
		yearEndDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
		if price, err := v.tickerManager.GetPrice(ctx, trades[0].Symbol, yearEndDate); err == nil {
			analysis.YearEndPosition = tax.Position{
				Date:     yearEndDate,
				Quantity: currentPosition,
				USDPrice: price,
			}
		} else {
			return analysis, common.NewHttpError("failed to get year end price", http.StatusInternalServerError)
		}
	}

	return analysis, nil
}
