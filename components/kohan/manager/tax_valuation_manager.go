package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

// TaxValuationManager handles currency exchange rate processing
type TaxValuationManager interface {
	ProcessValuations(ctx context.Context, valuations []tax.Valuation) ([]tax.TaxValuation, common.HttpError)
}

type TaxValuationManagerImpl struct {
	sbiManager SBIManager
}

// NewTaxValuationManager creates a new instance of TaxValuationManager
func NewTaxValuationManager(sbiManager SBIManager) *TaxValuationManagerImpl {
	return &TaxValuationManagerImpl{
		sbiManager: sbiManager,
	}
}

func (e *TaxValuationManagerImpl) ProcessValuations(ctx context.Context, valuations []tax.Valuation) ([]tax.TaxValuation, common.HttpError) {
	var taxValuations []tax.TaxValuation

	for _, valuation := range valuations {
		// Create base tax valuation
		taxValuation := tax.TaxValuation{
			Ticker: valuation.Ticker,
		}

		// Process each position with exchange rates
		if err := e.processPosition(ctx, &taxValuation.FirstPosition, valuation.FirstPosition); err != nil {
			return nil, err
		}
		if err := e.processPosition(ctx, &taxValuation.PeakPosition, valuation.PeakPosition); err != nil {
			return nil, err
		}
		if err := e.processPosition(ctx, &taxValuation.YearEndPosition, valuation.YearEndPosition); err != nil {
			return nil, err
		}

		taxValuations = append(taxValuations, taxValuation)
	}

	return taxValuations, nil
}

// processPosition converts a Position to TaxPosition by fetching and applying exchange rate
func (e *TaxValuationManagerImpl) processPosition(ctx context.Context, taxPosition *tax.TaxPosition, position tax.Position) common.HttpError {
	// Copy base position
	taxPosition.Position = position

	// Skip empty positions (e.g. when no year end position)
	if taxPosition.Quantity == 0 {
		return nil
	}

	// Set TTDate same as position date for rate lookup
	taxPosition.TTDate = position.Date

	// Get exchange rate for position date
	rate, err := e.sbiManager.GetTTBuyRate(taxPosition.TTDate)
	if err != nil {
		return err
	}

	// Set exchange rate for the position
	taxPosition.TTRate = rate

	return nil
}
