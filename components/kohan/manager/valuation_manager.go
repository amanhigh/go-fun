package manager

import (
	"context"
	"fmt"
	"time"

	"net/http"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type ValuationManager interface {
	AnalyzeValuation(ctx context.Context, trades []tax.Trade, year int) (tax.Valuation, error)
}

type ValuationManagerImpl struct {
	tickerManager TickerManager
}

func NewValuationManager(tickerManager TickerManager) ValuationManager {
	return &ValuationManagerImpl{
		tickerManager: tickerManager,
	}
}

func (v *ValuationManagerImpl) AnalyzeValuation(ctx context.Context, trades []tax.Trade, year int) (tax.Valuation, error) {
	if len(trades) == 0 {
		return tax.Valuation{}, common.NewHttpError("no trades provided", http.StatusBadRequest)
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
			firstBuyDate = t.Date
			analysis.FirstPosition = tax.Position{
				Date:     t.Date,
				Quantity: t.Quantity,
				USDPrice: t.USDPrice,
				USDValue: t.USDValue,
			}
		}

		// Track peak position
		if currentPosition > maxPosition {
			maxPosition = currentPosition
			analysis.PeakPosition = tax.Position{
				Date:     t.Date,
				Quantity: currentPosition,
				USDPrice: t.USDPrice,
				USDValue: currentPosition * t.USDPrice,
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
				USDValue: price * currentPosition,
			}
		} else {
			return analysis, fmt.Errorf("failed to get year end price: %w", err)
		}
	}

	return analysis, nil
}
