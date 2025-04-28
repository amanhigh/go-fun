package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type DividendManager interface {
	// Retrieves all Dividend records for the specified financial year.
	// The year parameter represents the starting year of the financial year (e.g., 2023 for FY 2023-24).
	GetDividendsForYear(ctx context.Context, year int) ([]tax.Dividend, common.HttpError) // Added method signature

	// Processes a list of Dividend records, adding INR values based on exchange rates.
	ProcessDividends(ctx context.Context, dividends []tax.Dividend) ([]tax.INRDividend, common.HttpError)
}

type DividendManagerImpl struct {
	exchangeManager      ExchangeManager
	financialYearManager FinancialYearManager[tax.Dividend] // Added
	dividendRepository   repository.DividendRepository      // Added
}

func NewDividendManager(
	exchangeManager ExchangeManager,
	financialYearManager FinancialYearManager[tax.Dividend], // Added
	dividendRepository repository.DividendRepository, // Added
) *DividendManagerImpl {
	return &DividendManagerImpl{
		exchangeManager:      exchangeManager,
		financialYearManager: financialYearManager, // Added
		dividendRepository:   dividendRepository,   // Added
	}
}

func (d *DividendManagerImpl) ProcessDividends(ctx context.Context, dividends []tax.Dividend) (inrDividends []tax.INRDividend, err common.HttpError) {
	exchangeables := make([]tax.Exchangeable, 0, len(dividends))
	inrDividends = make([]tax.INRDividend, len(dividends)) // Pre-allocate slice

	for i, dividend := range dividends {
		inrDividends[i].Dividend = dividend
		exchangeables = append(exchangeables, &inrDividends[i]) // Add pointer to element in pre-allocated slice
	}

	err = d.exchangeManager.Exchange(ctx, exchangeables)
	return
}

// GetDividendsForYear implementation added
func (d *DividendManagerImpl) GetDividendsForYear(ctx context.Context, year int) ([]tax.Dividend, common.HttpError) {
	// Get all records from repository
	records, err := d.dividendRepository.GetAllRecords(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by financial year
	return d.financialYearManager.FilterRecordsByFY(ctx, records, year)
}
