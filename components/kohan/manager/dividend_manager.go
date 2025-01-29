package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type DividendManager interface {
	ProcessDividends(ctx context.Context, dividends []tax.Dividend) ([]tax.INRDividend, common.HttpError)
}

type DividendManagerImpl struct {
	exchangeManager ExchangeManager
}

func NewDividendManager(exchangeManager ExchangeManager) *DividendManagerImpl {
	return &DividendManagerImpl{
		exchangeManager: exchangeManager,
	}
}

func (d *DividendManagerImpl) ProcessDividends(ctx context.Context, dividends []tax.Dividend) (inrDividends []tax.INRDividend, err common.HttpError) {
	exchangeables := make([]tax.Exchangeable, len(dividends))
	for _, dividend := range dividends {
		var inrDividend tax.INRDividend
		inrDividend.Dividend = dividend

		inrDividends = append(inrDividends, inrDividend)
		exchangeables = append(exchangeables, &inrDividend)
	}

	err = d.exchangeManager.Exchange(ctx, exchangeables)
	return
}
