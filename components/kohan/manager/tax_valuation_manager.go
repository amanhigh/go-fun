package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

// TaxValuationManager handles currency exchange rate processing
type TaxValuationManager interface {
	// Processes a list of Valuation records, adding INR values based on exchange rates.
	// Accepts dividends (can be empty slice) to calculate AmountPaid per ticker.
	ProcessValuations(ctx context.Context, valuations []tax.Valuation, dividends []tax.INRDividend) ([]tax.INRValuation, common.HttpError)

	// GetYearlyValuationsUSD calculates the base USD Valuation (First, Peak, YearEnd)
	// for all relevant tickers based on trade history up to the end of the specified calendar year.
	GetYearlyValuationsUSD(ctx context.Context, year int) ([]tax.Valuation, common.HttpError)
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

func (v *TaxValuationManagerImpl) ProcessValuations(ctx context.Context, valuations []tax.Valuation, dividends []tax.INRDividend) (inrValuations []tax.INRValuation, err common.HttpError) {
	// Pre-allocate slice with the correct size
	inrValuations = make([]tax.INRValuation, len(valuations))

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
	if err != nil {
		return
	}

	// Calculate AmountPaid (mandatory - will be 0 if no dividends)
	dividendsByTicker := v.groupDividendsByTicker(dividends)
	for i := range inrValuations {
		inrValuations[i].AmountPaid = dividendsByTicker[inrValuations[i].Ticker]
	}

	return
}

func (v *TaxValuationManagerImpl) groupDividendsByTicker(dividends []tax.INRDividend) map[string]float64 {
	sumByTicker := make(map[string]float64)
	for _, div := range dividends {
		sumByTicker[div.Symbol] += div.INRValue()
	}
	return sumByTicker
}

// GetYearlyValuationsUSD passes the call through to the underlying ValuationManager.
func (v *TaxValuationManagerImpl) GetYearlyValuationsUSD(ctx context.Context, year int) ([]tax.Valuation, common.HttpError) {
	return v.valuationManager.GetYearlyValuationsUSD(ctx, year)
}
