package manager

import (
	"context"
	"errors"

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

		requestedDate, dateErr := exchangeable.GetDate()
		if dateErr != nil {
			return dateErr
		}

		rate, err := e.sbiManager.GetTTBuyRate(ctx, requestedDate)

		var closestErr tax.ClosestDateError
		if err == nil {
			// Exact date found - No error
			exchangeable.SetTTRate(rate)
			exchangeable.SetTTDate(requestedDate)
		} else if errors.As(err, &closestErr) {
			// Handle closest date scenario specifically
			exchangeable.SetTTRate(rate) // Rate is still returned by GetTTBuyRate even with ClosestDateError
			exchangeable.SetTTDate(closestErr.GetClosestDate())
			log.Warn().Float64("RateSet", rate).Time("RequestedDate", requestedDate).Time("DateSet", closestErr.GetClosestDate()).Msg("ExchangeManager: Set closest rate/date")
		} else {
			// Handle any other non-nil error
			return err
		}
	}
	return nil
}
