package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/rs/zerolog/log"
)

//go:generate mockery --name ExchangeManager
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
		requestedDate := exchangeable.GetDate()
		rate, err := e.sbiManager.GetTTBuyRate(ctx, requestedDate)

		// Return early for non-closest-date errors
		if err != nil {
			if _, ok := err.(tax.ClosestDateError); !ok {
				return err
			}
		}

		// Handle closest date scenario
		if closestErr, ok := err.(tax.ClosestDateError); ok {
			log.Ctx(ctx).Warn().
				Time("RequestedDate", requestedDate).
				Time("ClosestDate", closestErr.GetClosestDate()).
				Msg("Using closest available exchange rate")

			exchangeable.SetTTRate(rate)
			exchangeable.SetTTDate(closestErr.GetClosestDate())
		} else {
			// Exact date found
			exchangeable.SetTTRate(rate)
			exchangeable.SetTTDate(requestedDate)
		}
	}
	return nil
}
