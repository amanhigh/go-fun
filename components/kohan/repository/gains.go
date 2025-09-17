package repository

import "github.com/amanhigh/go-fun/models/tax"

type GainsRepository interface {
	BaseCSVRepository[tax.Gains]
}

type GainsRepositoryImpl struct {
	*BaseCSVRepositoryImpl[tax.Gains]
}

func NewGainsRepository(filePath string) GainsRepository {
	return &GainsRepositoryImpl{
		BaseCSVRepositoryImpl: NewBaseCSVRepository[tax.Gains](filePath),
	}
}
