package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

// TaxValuationManager handles currency exchange rate processing
type TaxValuationManager interface {
	ProcessValuations(ctx context.Context, valuations []tax.Valuation) ([]tax.INRValutaion, common.HttpError)
}

type TaxValuationManagerImpl struct {
	exchangeManager ExchangeManager
}

// NewTaxValuationManager creates a new instance of TaxValuationManager
func NewTaxValuationManager(exchangeManager ExchangeManager) TaxValuationManager {
	return &TaxValuationManagerImpl{
		exchangeManager: exchangeManager,
	}
}

func (v *TaxValuationManagerImpl) ProcessValuations(ctx context.Context, valuations []tax.Valuation) ([]tax.INRValutaion, common.HttpError) {
	var result []tax.INRValutaion

	for _, valuation := range valuations {
		// Create tax valuation with positions
		taxValuation := tax.NewINRValuation(valuation)

		// Collect all positions that need exchange rates
		var positions []tax.Exchangeable
		if valuation.FirstPosition.Quantity > 0 {
			positions = append(positions, &taxValuation.FirstPosition)
		}
		if valuation.PeakPosition.Quantity > 0 {
			positions = append(positions, &taxValuation.PeakPosition)
		}
		if valuation.YearEndPosition.Quantity > 0 {
			positions = append(positions, &taxValuation.YearEndPosition)
		}

		// Process exchange rates for all positions
		if err := v.exchangeManager.Exchange(ctx, positions); err != nil {
			return nil, err
		}

		result = append(result, taxValuation)
	}

	return result, nil
}
