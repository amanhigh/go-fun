package manager

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/models/tax"
)

type TaxManager interface {
	ValueTickers(ctx context.Context, tickers []string, year int) ([]tax.Valuation, error)
}

type TaxManagerImpl struct {
	tickerManager TickerManager
	sbiManager    SBIManager
}

func NewTaxManager(tickerManager TickerManager, sbiManager SBIManager) *TaxManagerImpl {
	return &TaxManagerImpl{
		tickerManager: tickerManager,
		sbiManager:    sbiManager,
	}
}

func (f *TaxManagerImpl) ValueTickers(ctx context.Context, tickers []string, year int) ([]tax.Valuation, error) {
	var results []tax.Valuation
	for _, ticker := range tickers {
		// Get USD Analysis
		analysis, err := f.tickerManager.ValueTicker(ctx, ticker, year)
		if err != nil {
			return nil, err
		}

		// Get TT Rates
		// BUG: Use Date Constants
		peakDate, _ := time.Parse("2006-01-02", analysis.PeakDate)
		yearEndDate, _ := time.Parse("2006-01-02", analysis.YearEndDate)

		peakRate, err := f.sbiManager.GetTTBuyRate(peakDate)
		if err != nil {
			return nil, err
		}

		yearEndRate, err := f.sbiManager.GetTTBuyRate(yearEndDate)
		if err != nil {
			return nil, err
		}

		valuation := tax.Valuation{
			BaseValuation:   analysis,
			PeakTTRate:      peakRate,
			YearEndTTRate:   yearEndRate,
			PeakPriceINR:    analysis.PeakPrice * peakRate,
			YearEndPriceINR: analysis.YearEndPrice * yearEndRate,
		}

		results = append(results, valuation)
	}

	return results, nil
}
