package repository

import (
	"github.com/amanhigh/go-fun/models/tax"
)

type TradeRepository interface {
	BaseCSVRepository[tax.Trade]
}

type TradeRepositoryImpl struct {
	*BaseCSVRepositoryImpl[tax.Trade]
}

func NewTradeRepository(filePath string) TradeRepository {
	return &TradeRepositoryImpl{
		BaseCSVRepositoryImpl: NewBaseCSVRepository[tax.Trade](filePath),
	}
}
