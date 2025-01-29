package repository

import "github.com/amanhigh/go-fun/models/tax"

//go:generate mockery --name InterestRepository
type InterestRepository interface {
	BaseCSVRepository[tax.Interest]
}

type InterestRepositoryImpl struct {
	*BaseCSVRepositoryImpl[tax.Interest]
}

func NewInterestRepository(filePath string) InterestRepository {
	return &InterestRepositoryImpl{
		BaseCSVRepositoryImpl: NewBaseCSVRepository[tax.Interest](filePath),
	}
}
