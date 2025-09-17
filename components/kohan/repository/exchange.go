package repository

import "github.com/amanhigh/go-fun/models/tax"

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
