package manager

import (
	"context"
	"fmt"
	"net/http"

	repository "github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

//go:generate mockery --name AccountManager
type AccountManager interface {
	GetRecord(ctx context.Context, symbol string) (tax.Account, common.HttpError)
}

type AccountManagerImpl struct {
	repository repository.AccountRepository
}

func NewAccountManager(repo repository.AccountRepository) AccountManager {
	return &AccountManagerImpl{
		repository: repo,
	}
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
