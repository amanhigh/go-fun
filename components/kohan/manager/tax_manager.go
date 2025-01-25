package manager

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/models/fa"
)

type TaxManager interface {
	ProcessTickers(ctx context.Context, tickers []string, year int) ([]fa.TickerInfo, error)
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

func (f *TaxManagerImpl) ProcessTickers(ctx context.Context, tickers []string, year int) ([]fa.TickerInfo, error) {
	var results []fa.TickerInfo
	for _, ticker := range tickers {
		// Get USD Analysis
		analysis, err := f.tickerManager.AnalyzeTicker(ctx, ticker, year)
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

		// Create FATickerAnalysis
		faAnalysis := fa.TickerInfo{
			TickerAnalysis:  analysis,
			PeakTTRate:      peakRate,
			YearEndTTRate:   yearEndRate,
			PeakPriceINR:    analysis.PeakPrice * peakRate,
			YearEndPriceINR: analysis.YearEndPrice * yearEndRate,
		}

		results = append(results, faAnalysis)
	}

	return results, nil
}
