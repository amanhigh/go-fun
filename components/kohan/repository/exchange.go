package repository

import "github.com/amanhigh/go-fun/models/tax"

//go:generate mockery --name ExchangeRepository
type ExchangeRepository interface {
	BaseCSVRepository[tax.SbiRate]
}

type ExchangeRepositoryImpl struct {
	*BaseCSVRepositoryImpl[tax.SbiRate]
}

func NewExchangeRepository(filePath string) ExchangeRepository {
	return &ExchangeRepositoryImpl{
		BaseCSVRepositoryImpl: NewBaseCSVRepository[tax.SbiRate](filePath),
	}
}
