package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type ExchangeManager interface {
	Exchange(ctx context.Context, exchangeables []tax.Exchangeable) common.HttpError
}

type ExchangeManagerImpl struct {
	sbiManager SBIManager
}

func NewExchangeManager(sbiManager SBIManager) ExchangeManager {
	return &ExchangeManagerImpl{
		sbiManager: sbiManager,
	}
}

func (e *ExchangeManagerImpl) Exchange(ctx context.Context, exchangeables []tax.Exchangeable) common.HttpError {
	for _, exchangeable := range exchangeables {
		date := exchangeable.GetDate()

		// Get exchange rate for date
		rate, err := e.sbiManager.GetTTBuyRate(date)
		if err != nil {
			return err
		}

		// Set exchange rate and date
		exchangeable.SetTTRate(rate)
		exchangeable.SetTTDate(date)
	}
	return nil
}
