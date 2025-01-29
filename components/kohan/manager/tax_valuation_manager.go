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

func (v *TaxValuationManagerImpl) ProcessValuations(ctx context.Context, valuations []tax.Valuation) (inrValuations []tax.INRValutaion, err common.HttpError) {
	exchangeAbles := make([]tax.Exchangeable, 0, len(valuations))

	for _, valuation := range valuations {
		// Create tax valuation with positions
		inrValuation := tax.NewINRValuation(valuation)

		// Collect all positions that need exchange rates
		exchangeAbles = append(exchangeAbles, &inrValuation.FirstPosition)
		exchangeAbles = append(exchangeAbles, &inrValuation.PeakPosition)
		exchangeAbles = append(exchangeAbles, &inrValuation.YearEndPosition)

		inrValuations = append(inrValuations, inrValuation)
	}

	err = v.exchangeManager.Exchange(ctx, exchangeAbles)

	return
}
