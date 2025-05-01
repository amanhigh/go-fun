package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type TaxManager interface {
	GetTaxSummary(ctx context.Context, year int) (tax.Summary, common.HttpError)
}

// Struct updated with TaxValuationManager
type TaxManagerImpl struct {
	capitalGainManager  CapitalGainManager
	dividendManager     DividendManager
	interestManager     InterestManager
	taxValuationManager TaxValuationManager // Added field
}

// Constructor updated to accept TaxValuationManager
func NewTaxManager(
	capitalGainManager CapitalGainManager,
	dividendManager DividendManager,
	interestManager InterestManager,
	taxValuationManager TaxValuationManager, // Added parameter
) TaxManager {
	return &TaxManagerImpl{
		capitalGainManager:  capitalGainManager,
		dividendManager:     dividendManager,
		interestManager:     interestManager,
		taxValuationManager: taxValuationManager, // Assign new dependency
	}
}

func (t *TaxManagerImpl) GetTaxSummary(ctx context.Context, year int) (summary tax.Summary, err common.HttpError) {
	// Process gains
	if summary.INRGains, err = t.processGains(ctx, year); err != nil {
		return
	}

	// Process dividends
	if summary.INRDividends, err = t.processDividends(ctx, year); err != nil {
		return
	}

	// Process interest
	if summary.INRInterest, err = t.processInterest(ctx, year); err != nil {
		return
	}

	// Process valuations
	if summary.INRValuations, err = t.processValuations(ctx, year); err != nil {
		return
	}

	return summary, nil
}

func (t *TaxManagerImpl) processGains(ctx context.Context, year int) ([]tax.INRGains, common.HttpError) {
	gains, err := t.capitalGainManager.GetGainsForYear(ctx, year)
	if err != nil {
		return nil, err
	}
	return t.capitalGainManager.ProcessTaxGains(ctx, gains)
}

func (t *TaxManagerImpl) processDividends(ctx context.Context, year int) ([]tax.INRDividend, common.HttpError) {
	dividends, err := t.dividendManager.GetDividendsForYear(ctx, year)
	if err != nil {
		return nil, err
	}
	return t.dividendManager.ProcessDividends(ctx, dividends)
}

func (t *TaxManagerImpl) processInterest(ctx context.Context, year int) ([]tax.INRInterest, common.HttpError) {
	interests, err := t.interestManager.GetInterestForYear(ctx, year)
	if err != nil {
		return nil, err
	}
	return t.interestManager.ProcessInterest(ctx, interests)
}

func (t *TaxManagerImpl) processValuations(ctx context.Context, year int) ([]tax.INRValutaion, common.HttpError) {
	usdValuations, err := t.taxValuationManager.GetYearlyValuationsUSD(ctx, year)
	if err != nil {
		return nil, err
	}
	return t.taxValuationManager.ProcessValuations(ctx, usdValuations)
}
