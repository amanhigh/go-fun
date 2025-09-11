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

func (a *AccountManagerImpl) GenerateYearEndAccounts(_ context.Context, year int, valuations []tax.Valuation) common.HttpError {
	accounts := tax.FromValuations(valuations)

	// Create a new file for the year-end accounts
	fileName := fmt.Sprintf("accounts_%d.csv", year)
	filePath := filepath.Join(a.accountDir, fileName)

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

func (a *AccountManagerImpl) GetRecord(ctx context.Context, symbol string, year int) (account tax.Account, err common.HttpError) {
	// Smart detection: Only check for auto-generated previous year file
	prevYearPath := fmt.Sprintf("accounts_%d.csv", year-1)
	fallbackPath := filepath.Join(a.accountDir, prevYearPath)

	if _, fileErr := os.Stat(fallbackPath); fileErr == nil {
		// Use auto-generated previous year file
		prevRepo := repository.NewAccountRepository(fallbackPath)
		records, repoErr := prevRepo.GetRecordsForTicker(ctx, symbol)
		if repoErr == nil && len(records) > 0 {
			// Validate single record from previous year file
			switch len(records) {
			case 1:
				return records[0], nil
			default:
				return account, common.NewHttpError(fmt.Sprintf("multiple accounts found for symbol: %s in %s", symbol, prevYearPath), http.StatusBadRequest)
			}
		}
	}

	// No previous year file OR ticker not found -> Fresh start
	return tax.Account{}, common.ErrNotFound
}
