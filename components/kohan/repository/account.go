package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
)

//go:generate mockery --name=AccountRepository
type AccountRepository interface {
	GetAllRecordsForYear(ctx context.Context, year int) ([]tax.Account, common.HttpError)
	SaveYearEndAccounts(ctx context.Context, year int, accounts []tax.Account) common.HttpError
}

type AccountRepositoryImpl struct {
	accountDir string
}

func NewAccountRepository(accountDir string) AccountRepository {
	return &AccountRepositoryImpl{
		accountDir: accountDir,
	}
}

func (r *AccountRepositoryImpl) GetAllRecordsForYear(ctx context.Context, year int) ([]tax.Account, common.HttpError) {
	prevYearPath := filepath.Join(r.accountDir, fmt.Sprintf("accounts_%d.csv", year-1))

	// File existence check
	if _, err := os.Stat(prevYearPath); err != nil {
		return nil, common.ErrNotFound
	}

	// Read accounts from file
	return r.readAccountsFromFile(ctx, prevYearPath)
}

func (r *AccountRepositoryImpl) SaveYearEndAccounts(_ context.Context, year int, accounts []tax.Account) common.HttpError {
	fileName := fmt.Sprintf("accounts_%d.csv", year)
	filePath := filepath.Join(r.accountDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return common.NewServerError(fmt.Errorf("failed to create year-end accounts file: %w", err))
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&accounts, file); err != nil {
		return common.NewServerError(fmt.Errorf("failed to write year-end accounts: %w", err))
	}

	return nil
}

// Private method with only the CSV logic we need
func (r *AccountRepositoryImpl) readAccountsFromFile(ctx context.Context, filePath string) ([]tax.Account, common.HttpError) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("path", filePath).Msg("Failed to open accounts CSV")
		return nil, common.NewServerError(err)
	}
	defer file.Close()

	var accounts []tax.Account
	if err := gocsv.UnmarshalFile(file, &accounts); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse accounts CSV")
		return nil, common.NewServerError(err)
	}

	if len(accounts) == 0 {
		return nil, common.ErrNotFound // Empty file = no records for this year
	}

	return accounts, nil
}
