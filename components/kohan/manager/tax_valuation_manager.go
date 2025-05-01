package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

// TaxValuationManager handles currency exchange rate processing
type TaxValuationManager interface {
	// Processes a list of Valuation records, adding INR values based on exchange rates.
	ProcessValuations(ctx context.Context, valuations []tax.Valuation) ([]tax.INRValutaion, common.HttpError)

	// GetYearlyValuationsUSD calculates the base USD Valuation (First, Peak, YearEnd)
	// for all relevant tickers based on trade history up to the end of the specified calendar year.
	GetYearlyValuationsUSD(ctx context.Context, year int) ([]tax.Valuation, common.HttpError) // Added method signature
}

// Implementation struct updated with ValuationManager
type TaxValuationManagerImpl struct {
	exchangeManager  ExchangeManager
	valuationManager ValuationManager
}

// NewTaxValuationManager creates a new instance of TaxValuationManager
func NewTaxValuationManager(exchangeManager ExchangeManager, valuationManager ValuationManager) TaxValuationManager {
	return &TaxValuationManagerImpl{
		exchangeManager:  exchangeManager,
		valuationManager: valuationManager,
	}
}

func (v *TaxValuationManagerImpl) ProcessValuations(ctx context.Context, valuations []tax.Valuation) (inrValuations []tax.INRValutaion, err common.HttpError) {
	exchangeAbles := make([]tax.Exchangeable, 0, len(valuations))

	for _, valuation := range valuations {
		// Create tax valuation with positions
		inrValuation := tax.NewINRValuation(valuation)

		// Collect all positions that need exchange rates
		exchangeAbles = append(exchangeAbles,
			&inrValuation.FirstPosition,
			&inrValuation.PeakPosition,
			&inrValuation.YearEndPosition)

		inrValuations = append(inrValuations, inrValuation)
	}

	err = v.exchangeManager.Exchange(ctx, exchangeAbles)

	return
}

// GetYearlyValuationsUSD passes the call through to the underlying ValuationManager.
func (v *TaxValuationManagerImpl) GetYearlyValuationsUSD(ctx context.Context, year int) ([]tax.Valuation, common.HttpError) {
	return v.valuationManager.GetYearlyValuationsUSD(ctx, year)
}
