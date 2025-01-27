package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

// ExchangeManager handles currency exchange rate processing
type ExchangeManager interface {
	ProcessValuations(ctx context.Context, valuations []tax.Valuation) ([]tax.TaxValuation, common.HttpError)
}

type ExchangeManagerImpl struct {
	sbiManager SBIManager
}

// NewExchangeManager creates a new instance of ExchangeManager
func NewExchangeManager(sbiManager SBIManager) *ExchangeManagerImpl {
	return &ExchangeManagerImpl{
		sbiManager: sbiManager,
	}
}

func (e *ExchangeManagerImpl) ProcessValuations(ctx context.Context, valuations []tax.Valuation) ([]tax.TaxValuation, common.HttpError) {
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
func (e *ExchangeManagerImpl) processPosition(ctx context.Context, taxPosition *tax.TaxPosition, position tax.Position) common.HttpError {
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
	taxPosition.TTBuyRate = rate

	return nil
}
