package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type TaxManager interface {
	GetTaxSummary(ctx context.Context, year int) (tax.Summary, common.HttpError)
}

// Struct updated with InterestManager
type TaxManagerImpl struct {
	capitalGainManager CapitalGainManager
	dividendManager    DividendManager
	interestManager    InterestManager // Added field
}

// Constructor updated to accept InterestManager
func NewTaxManager(
	capitalGainManager CapitalGainManager,
	dividendManager DividendManager,
	interestManager InterestManager, // Added parameter
) TaxManager {
	return &TaxManagerImpl{
		capitalGainManager: capitalGainManager,
		dividendManager:    dividendManager,
		interestManager:    interestManager, // Assign new dependency
	}
}

func (t *TaxManagerImpl) GetTaxSummary(ctx context.Context, year int) (summary tax.Summary, err common.HttpError) {
	// Get and process gains for the year
	gains, err := t.capitalGainManager.GetGainsForYear(ctx, year)
	if err != nil {
		return summary, err
	}

	// Process gains to INR
	summary.INRGains, err = t.capitalGainManager.ProcessTaxGains(ctx, gains)
	if err != nil {
		return summary, err
	}

	// Get and process dividends for the year
	dividends, err := t.dividendManager.GetDividendsForYear(ctx, year)
	if err != nil {
		return summary, err
	}
	summary.INRDividends, err = t.dividendManager.ProcessDividends(ctx, dividends)
	if err != nil {
		return summary, err
	}

	// Get and process interest for the specific year (NEW SECTION)
	interests, err := t.interestManager.GetInterestForYear(ctx, year)
	if err != nil {
		return summary, err
	}
	summary.INRInterest, err = t.interestManager.ProcessInterest(ctx, interests)
	if err != nil {
		return summary, err
	}

	// Assign other processed data (Valuation etc. if they exist) to summary here...
	// summary.INRPositions = ...

	return summary, nil
}
