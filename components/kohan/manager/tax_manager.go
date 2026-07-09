package manager

import (
	"context"
	"time"

	"fmt"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type TaxManager interface {
	GetTaxSummary(ctx context.Context, year int) (tax.Summary, common.HttpError)
	SaveTaxSummaryToExcel(ctx context.Context, year int, summary tax.Summary) error
}

type TaxManagerImpl struct {
	capitalGainManager  CapitalGainManager
	dividendManager     DividendManager
	interestManager     InterestManager
	taxValuationManager TaxValuationManager
	excelManager        ExcelManager
	accountManager      AccountManager
	sbiManager          SBIManager
}

//nolint:revive // argument-limit: 7 params matches existing pattern
func NewTaxManager(
	capitalGainManager CapitalGainManager,
	dividendManager DividendManager,
	interestManager InterestManager,
	taxValuationManager TaxValuationManager,
	excelManager ExcelManager,
	accountManager AccountManager,
	sbiManager SBIManager,
) TaxManager {
	return &TaxManagerImpl{
		capitalGainManager:  capitalGainManager,
		dividendManager:     dividendManager,
		interestManager:     interestManager,
		taxValuationManager: taxValuationManager,
		excelManager:        excelManager,
		accountManager:      accountManager,
		sbiManager:          sbiManager,
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

	// Process SBI TT month-end rates for the FY (Apr→Mar)
	if summary.TTMonthEndRates, err = t.processTTMonthEndRates(ctx, year); err != nil {
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

func (t *TaxManagerImpl) processValuations(ctx context.Context, year int) ([]tax.INRValuation, common.HttpError) {
	usdValuations, err := t.taxValuationManager.GetYearlyValuationsUSD(ctx, year)
	if err != nil {
		return nil, err
	}

	// Generate Year End Accounts CSV
	if accountErr := t.accountManager.GenerateYearEndAccounts(ctx, year, usdValuations); accountErr != nil {
		return nil, accountErr
	}

	// Get US year dividends for AmountPaid calculation (always fetch, even if empty)
	usDividends, err := t.dividendManager.GetDividendsForUSYear(ctx, year)
	if err != nil {
		return nil, err
	}

	// Process dividends to INR (always pass, even if empty slice)
	inrDividends, err := t.dividendManager.ProcessDividends(ctx, usDividends)
	if err != nil {
		return nil, err
	}

	// Always pass dividends to ProcessValuations (mandatory parameter)
	return t.taxValuationManager.ProcessValuations(ctx, usdValuations, inrDividends)
}

// processTTMonthEndRates populates one SBI TT Buy month-end rate for each month
// of the Indian financial year (April of "year" to March of "year+1").
// Returns exactly 12 rates in Apr→Mar order by reusing the existing GetLastMonthEndRate method.
func (t *TaxManagerImpl) processTTMonthEndRates(ctx context.Context, year int) ([]tax.MonthEndRate, common.HttpError) {
	rates := make([]tax.MonthEndRate, 0, 12)

	for i := range 12 {
		monthDate := time.Date(year, time.April, 1, 0, 0, 0, 0, time.UTC).AddDate(0, i, 0)

		rate, err := t.sbiManager.GetLastMonthEndRate(ctx, monthDate)
		if err != nil {
			return nil, err
		}

		rates = append(rates, rate)
	}

	return rates, nil
}

func (t *TaxManagerImpl) SaveTaxSummaryToExcel(ctx context.Context, year int, summary tax.Summary) error {
	if err := t.excelManager.GenerateTaxSummaryExcel(ctx, year, summary); err != nil {
		return fmt.Errorf("failed to generate tax summary excel: %w", err)
	}
	return nil
}
