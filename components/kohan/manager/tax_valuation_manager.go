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
	// Pre-allocate slice with the correct size
	inrValuations = make([]tax.INRValutaion, len(valuations))

	// Pre-allocate exchangeAbles slice with the final capacity
	exchangeAbles := make([]tax.Exchangeable, 0, len(valuations)*3) // 3 positions per valuation

	// Iterate using index
	for i, valuation := range valuations {
		// Create and assign the INRValuation directly into the final slice at index i
		inrValuations[i] = tax.NewINRValuation(valuation)

		// Append pointers *from the element within the inrValuations slice*
		exchangeAbles = append(exchangeAbles,
			&inrValuations[i].FirstPosition,
			&inrValuations[i].PeakPosition,
			&inrValuations[i].YearEndPosition)
	}

	// ExchangeManager modifies the structs in inrValuations via these pointers
	err = v.exchangeManager.Exchange(ctx, exchangeAbles)
	return
}

// GetYearlyValuationsUSD passes the call through to the underlying ValuationManager.
func (v *TaxValuationManagerImpl) GetYearlyValuationsUSD(ctx context.Context, year int) ([]tax.Valuation, common.HttpError) {
	return v.valuationManager.GetYearlyValuationsUSD(ctx, year)
}
