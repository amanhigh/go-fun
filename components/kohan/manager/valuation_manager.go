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

	if err := v.validateTrades(trades); err != nil {
		return tax.Valuation{}, err
	}

	analysis, currentPosition := v.processTradesAndUpdateAnalysis(trades)
	analysis, err := v.updateYearEndPosition(ctx, analysis, currentPosition, year)
	if err != nil {
		return analysis, err
	}

	return analysis, nil
}

func (v *ValuationManagerImpl) validateTrades(trades []tax.Trade) common.HttpError {
	for _, t := range trades {
		if t.Symbol != trades[0].Symbol {
			return common.NewHttpError("multiple tickers found in trades", http.StatusBadRequest)
		}
	}
	return nil
}

func (v *ValuationManagerImpl) processTradesAndUpdateAnalysis(trades []tax.Trade) (tax.Valuation, float64) {
	analysis := tax.Valuation{
		Ticker: trades[0].Symbol,
	}

	var currentPosition float64
	var maxPosition float64
	var firstBuyDate time.Time

	for _, t := range trades {
		if t.Type == "BUY" {
			currentPosition += t.Quantity
		} else {
			currentPosition -= t.Quantity
		}

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

		if currentPosition > maxPosition {
			maxPosition = currentPosition
			if parsedTime, err := t.GetDate(); err == nil {
				analysis.PeakPosition = tax.Position{
					Date:     parsedTime,
					Quantity: maxPosition,
					USDPrice: t.USDPrice,
				}
			}
		}
	}

	return analysis, currentPosition
}

func (v *ValuationManagerImpl) updateYearEndPosition(ctx context.Context, analysis tax.Valuation, currentPosition float64, year int) (tax.Valuation, common.HttpError) {
	if currentPosition > 0 {
		yearEndDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
		price, err := v.tickerManager.GetPrice(ctx, analysis.Ticker, yearEndDate)
		if err != nil {
			return analysis, common.NewHttpError("failed to get year end price", http.StatusInternalServerError)
		}
		analysis.YearEndPosition = tax.Position{
			Date:     yearEndDate,
			Quantity: currentPosition,
			USDPrice: price,
		}
	}
	return analysis, nil
}
