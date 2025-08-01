package manager

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	repository "github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/gocarina/gocsv"
)

//go:generate mockery --name AccountManager
type AccountManager interface {
	GetRecord(ctx context.Context, symbol string) (tax.Account, common.HttpError)
	GenerateYearEndAccounts(ctx context.Context, year int, valuations []tax.Valuation) common.HttpError
}

type AccountManagerImpl struct {
	repository      repository.AccountRepository
	accountFilePath string
}

func NewAccountManager(repo repository.AccountRepository, accountFilePath string) AccountManager {
	return &AccountManagerImpl{
		repository:      repo,
		accountFilePath: accountFilePath,
	}
}

func (a *AccountManagerImpl) GenerateYearEndAccounts(_ context.Context, year int, valuations []tax.Valuation) common.HttpError {
	accounts := tax.FromValuations(valuations)

	// Create a new file for the year-end accounts
	fileName := fmt.Sprintf("accounts_%d.csv", year)
	filePath := filepath.Join(filepath.Dir(a.accountFilePath), fileName)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return common.NewServerError(fmt.Errorf("failed to create year-end accounts file: %w", err))
	}
	defer file.Close()

	// Marshal and write to the new file
	if err := gocsv.MarshalFile(&accounts, file); err != nil {
		return common.NewServerError(fmt.Errorf("failed to write year-end accounts: %w", err))
	}

	return nil
}

func (a *AccountManagerImpl) GetRecord(ctx context.Context, symbol string) (account tax.Account, err common.HttpError) {
	// Get all records for symbol
	records, err := a.repository.GetRecordsForTicker(ctx, symbol)
	if err != nil {
		return account, err
	}

	// Validate single record
	switch len(records) {
	case 0:
		return account, common.ErrNotFound
	case 1:
		return records[0], nil
	default:
		return account, common.NewHttpError(fmt.Sprintf("multiple accounts found for symbol: %s", symbol), http.StatusBadRequest)
	}
}
