package repository

import (
	"github.com/amanhigh/go-fun/models/tax"
)

//go:generate mockery --name AccountRepository
type AccountRepository interface {
	BaseCSVRepository[tax.Account]
}

type AccountRepositoryImpl struct {
	*BaseCSVRepositoryImpl[tax.Account]
}

func NewAccountRepository(filePath string) AccountRepository {
	return &AccountRepositoryImpl{
		BaseCSVRepositoryImpl: NewBaseCSVRepository[tax.Account](filePath),
	}
}
