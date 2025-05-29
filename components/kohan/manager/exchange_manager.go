package manager

import (
	"context"
	"errors"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/rs/zerolog/log"
)

//go:generate mockery --name ExchangeManager
type ExchangeManager interface {
	Exchange(ctx context.Context, exchangeables []tax.Exchangeable) common.HttpError

	// ExchangeGains applies exchange rates to INRGains items.
	// For each gain, it uses the SBI TT Buy Rate from the last day of the month
	// immediately preceding the month of that gain's SellDate.
	// If an exact rate for that date is not found, it relies on the underlying
	// SBI rate provider to furnish a rate for the closest available date.
	ExchangeGains(ctx context.Context, gains []tax.INRGains) common.HttpError
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
		switch {
		case err == nil:
			// Exact date found - No error
			exchangeable.SetTTRate(rate)
			exchangeable.SetTTDate(requestedDate)
		case errors.As(err, &closestErr):
			// Handle closest date scenario specifically
			exchangeable.SetTTRate(rate) // Rate is still returned by GetTTBuyRate even with ClosestDateError
			exchangeable.SetTTDate(closestErr.GetClosestDate())
			log.Warn().Float64("RateSet", rate).Time("RequestedDate", requestedDate).Time("DateSet", closestErr.GetClosestDate()).Msg("ExchangeManager: Set closest rate/date")
		default:
			// Handle any other non-nil error
			return err
		}
	}
	return nil
}

// getExchangeRateTargetDateForGain calculates the target date for fetching
// an exchange rate for a capital gain, which is the last day of the
// month immediately preceding the month of the gain's sell date.
func getExchangeRateTargetDateForGain(sellDate time.Time) time.Time {
	firstDayOfSellMonth := time.Date(sellDate.Year(), sellDate.Month(), 1, 0, 0, 0, 0, sellDate.Location())
	lastDayOfPrecedingMonth := firstDayOfSellMonth.AddDate(0, 0, -1)
	return lastDayOfPrecedingMonth
}

func (e *ExchangeManagerImpl) ExchangeGains(ctx context.Context, gains []tax.INRGains) common.HttpError {
	for i := range gains {
		// Operate on the pointer to modify the original slice element
		gain := &gains[i]

		parsedSellDate, dateErr := gain.ParseSellDate()
		if dateErr != nil {
			return tax.NewInvalidDateError(dateErr.Error())
		}

		targetRateDate := getExchangeRateTargetDateForGain(parsedSellDate)

		rate, err := e.sbiManager.GetTTBuyRate(ctx, targetRateDate)

		var closestErr tax.ClosestDateError
		switch {
		case err == nil:
			// Exact date found - No error
			gain.TTRate = rate
			gain.TTDate = targetRateDate
		case errors.As(err, &closestErr):
			// Handle closest date scenario specifically
			gain.TTRate = rate // Rate is still returned by GetTTBuyRate even with ClosestDateError
			gain.TTDate = closestErr.GetClosestDate()
			log.Warn().
				Float64("RateSet", rate).
				Time("RequestedDate", targetRateDate).
				Time("DateSet", closestErr.GetClosestDate()).
				Str("Symbol", gain.Symbol).
				Msg("ExchangeManager.ExchangeGains: Set closest rate/date for gain")
		default:
			// So, if err is not nil and not ClosestDateError, it must be another common.HttpError from sbiManager.
			return err
		}
	}
	return nil
}
