package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type DividendManager interface {
	// Retrieves all Dividend records for the specified Indian financial year (Apr-Mar).
	// The year parameter represents the starting year of the financial year (e.g., 2023 for FY 2023-24).
	GetDividendsForYear(ctx context.Context, year int) ([]tax.Dividend, common.HttpError)

	// Retrieves all Dividend records for the specified US calendar year (Jan-Dec).
	// Used for Schedule FA reporting which follows US calendar year.
	GetDividendsForUSYear(ctx context.Context, year int) ([]tax.Dividend, common.HttpError)

	// Processes a list of Dividend records, adding INR values based on exchange rates.
	ProcessDividends(ctx context.Context, dividends []tax.Dividend) ([]tax.INRDividend, common.HttpError)
}

type DividendManagerImpl struct {
	exchangeManager      ExchangeManager
	financialYearManager FinancialYearManager[tax.Dividend]
	dividendRepository   repository.DividendRepository
}

func NewDividendManager(
	exchangeManager ExchangeManager,
	financialYearManager FinancialYearManager[tax.Dividend],
	dividendRepository repository.DividendRepository,
) *DividendManagerImpl {
	return &DividendManagerImpl{
		exchangeManager:      exchangeManager,
		financialYearManager: financialYearManager,
		dividendRepository:   dividendRepository,
	}
}

func (d *DividendManagerImpl) ProcessDividends(ctx context.Context, dividends []tax.Dividend) (inrDividends []tax.INRDividend, err common.HttpError) {
	exchangeables := make([]tax.Exchangeable, 0, len(dividends))
	inrDividends = make([]tax.INRDividend, len(dividends))

	for i, dividend := range dividends {
		inrDividends[i].Dividend = dividend
		exchangeables = append(exchangeables, &inrDividends[i])
	}

	err = d.exchangeManager.ExchangeWithPrecedingMonth(ctx, exchangeables)
	return
}

func (d *DividendManagerImpl) GetDividendsForYear(ctx context.Context, year int) ([]tax.Dividend, common.HttpError) {
	// Get all records from repository
	records, err := d.dividendRepository.GetAllRecords(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by Indian financial year (Apr-Mar)
	return d.financialYearManager.FilterIndia(ctx, records, year)
}

func (d *DividendManagerImpl) GetDividendsForUSYear(ctx context.Context, year int) ([]tax.Dividend, common.HttpError) {
	// Get all records from repository
	records, err := d.dividendRepository.GetAllRecords(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by US calendar year (Jan-Dec)
	return d.financialYearManager.FilterUS(ctx, records, year)
}
