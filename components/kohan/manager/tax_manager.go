package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type TaxManager interface {
	GetTaxSummary(ctx context.Context, year int) (tax.Summary, common.HttpError)
}

type TaxManagerImpl struct {
	capitalGainManager CapitalGainManager
	dividendManager    DividendManager // Added field
}

func NewTaxManager(capitalGainManager CapitalGainManager, dividendManager DividendManager) TaxManager {
	return &TaxManagerImpl{
		capitalGainManager: capitalGainManager,
		dividendManager:    dividendManager, // Added assignment
	}
}

func (t *TaxManagerImpl) GetTaxSummary(ctx context.Context, year int) (summary tax.Summary, err common.HttpError) {
	// Get and process gains for the year
	gains, err := t.capitalGainManager.GetGainsForYear(ctx, year)
	if err != nil {
		return summary, err
	}

	// Process gains to INR
	inrGains, err := t.capitalGainManager.ProcessTaxGains(ctx, gains)
	if err != nil {
		return summary, err
	}

	// Build summary
	summary.INRGains = inrGains

	// Get Dividends for the year (Added)
	dividends, httpErr := t.dividendManager.GetDividendsForYear(ctx, year)
	if httpErr != nil {
		return summary, httpErr // Return HttpError directly
	}

	// Process Dividends (Added)
	summary.INRDividends, httpErr = t.dividendManager.ProcessDividends(ctx, dividends) // Pass fetched dividends
	if httpErr != nil {
		// Log error? Fail on first error.
		return summary, httpErr // Return HttpError directly
	}

	return summary, nil
}
