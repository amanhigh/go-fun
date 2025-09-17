package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	repository "github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type AccountManager interface {
	GetRecord(ctx context.Context, symbol string, year int) (tax.Account, common.HttpError)
	GenerateYearEndAccounts(ctx context.Context, year int, valuations []tax.Valuation) common.HttpError
}

type AccountManagerImpl struct {
	repository repository.AccountRepository
	accountDir string
}

func NewAccountManager(repo repository.AccountRepository, accountDir string) AccountManager {
	return &AccountManagerImpl{
		repository: repo,
		accountDir: accountDir,
	}
}

func (a *AccountManagerImpl) GenerateYearEndAccounts(ctx context.Context, year int, valuations []tax.Valuation) common.HttpError {
	// Business logic: convert valuations to accounts
	accounts := tax.FromValuations(valuations)

	// Delegate write operation to repository (repository handles accountDir internally)
	return a.repository.SaveYearEndAccounts(ctx, year, accounts)
}

func (a *AccountManagerImpl) GetRecord(ctx context.Context, symbol string, year int) (tax.Account, common.HttpError) {
	// Get ALL records from repository for the specified year
	allRecords, err := a.repository.GetAllRecordsForYear(ctx, year)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return tax.Account{}, common.ErrNotFound // Fresh start
		}
		return tax.Account{}, err
	}

	// MANAGER responsibility: Filter by symbol
	var matchingRecords []tax.Account
	for _, record := range allRecords {
		if record.Symbol == symbol {
			matchingRecords = append(matchingRecords, record)
		}
	}

	// MANAGER responsibility: Business validation
	switch len(matchingRecords) {
	case 0:
		return tax.Account{}, common.ErrNotFound
	case 1:
		return matchingRecords[0], nil
	default:
		return tax.Account{}, common.NewHttpError(fmt.Sprintf("multiple accounts found for symbol: %s", symbol), http.StatusBadRequest)
	}
}
