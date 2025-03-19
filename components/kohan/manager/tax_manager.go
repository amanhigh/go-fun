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
}

func NewTaxManager(capitalGainManager CapitalGainManager) TaxManager {
	return &TaxManagerImpl{
		capitalGainManager: capitalGainManager,
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
	return summary, nil
}
