package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/repository" // Ensure repository import exists
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

// Interface definition updated
type InterestManager interface {
	// Retrieves all Interest records for the specified financial year.
	// The year parameter represents the starting year of the financial year (e.g., 2023 for FY 2023-24).
	GetInterestForYear(ctx context.Context, year int) ([]tax.Interest, common.HttpError)

	// Processes a list of Interest records, adding INR values based on exchange rates.
	ProcessInterest(ctx context.Context, interest []tax.Interest) ([]tax.INRInterest, common.HttpError)
}

// Implementation struct updated with new dependencies
type InterestManagerImpl struct {
	exchangeManager      ExchangeManager
	financialYearManager FinancialYearManager[tax.Interest]
	interestRepository   repository.InterestRepository
}

// Constructor updated to accept new dependencies
func NewInterestManager(
	exchangeManager ExchangeManager,
	financialYearManager FinancialYearManager[tax.Interest],
	interestRepository repository.InterestRepository,
) *InterestManagerImpl {
	return &InterestManagerImpl{
		exchangeManager:      exchangeManager,
		financialYearManager: financialYearManager,
		interestRepository:   interestRepository,
	}
}

// GetInterestForYear implementation added
func (i *InterestManagerImpl) GetInterestForYear(ctx context.Context, year int) ([]tax.Interest, common.HttpError) {
	// Get all records from repository
	records, err := i.interestRepository.GetAllRecords(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by financial year
	return i.financialYearManager.FilterIndia(ctx, records, year)
}

// ProcessInterest implementation (Ensure it uses injected exchangeManager)
func (i *InterestManagerImpl) ProcessInterest(ctx context.Context, interests []tax.Interest) (inrInterests []tax.INRInterest, err common.HttpError) {
	exchangeables := make([]tax.Exchangeable, 0, len(interests))
	inrInterests = make([]tax.INRInterest, len(interests)) // Pre-allocate slice

	for idx, interest := range interests {
		inrInterests[idx].Interest = interest
		exchangeables = append(exchangeables, &inrInterests[idx]) // Add pointer to element in pre-allocated slice
	}

	err = i.exchangeManager.Exchange(ctx, exchangeables)
	return
}
