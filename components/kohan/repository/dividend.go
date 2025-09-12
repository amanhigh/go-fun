package repository

import "github.com/amanhigh/go-fun/models/tax"

type DividendRepository interface {
	BaseCSVRepository[tax.Dividend]
}

type DividendRepositoryImpl struct {
	*BaseCSVRepositoryImpl[tax.Dividend]
}

func NewDividendRepository(filePath string) DividendRepository {
	return &DividendRepositoryImpl{
		BaseCSVRepositoryImpl: NewBaseCSVRepository[tax.Dividend](filePath),
	}
}
