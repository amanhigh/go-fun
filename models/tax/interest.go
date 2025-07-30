package tax

import (
	"fmt"
	"math"
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

type Interest struct {
	Symbol string  `csv:"Symbol"`
	Date   string  `csv:"Date"`
	Amount float64 `csv:"Amount"`
	Tax    float64 `csv:"Tax"`
	Net    float64 `csv:"Net"`
}

// CSVRecord interface implementation
func (i Interest) GetKey() string {
	return i.Symbol
}

func (i Interest) IsValid() bool {
	return i.Symbol != "" && i.Date != "" && i.Amount != 0
}

func (i Interest) GetDate() (time.Time, common.HttpError) {
	t, err := time.Parse(time.DateOnly, i.Date)
	if err != nil {
		return time.Time{}, NewInvalidDateError(fmt.Sprintf("failed to parse date '%s': %v", i.Date, err))
	}
	return t, nil
}

// INRInterest adds exchange rate details to basic interest
type INRInterest struct {
	Interest           // Embed original interest
	TTDate   time.Time // Date for exchange rate
	TTRate   float64   // Exchange rate
}

// Implement Exchangeable interface
func (i *INRInterest) GetDate() (time.Time, common.HttpError) {
	// Use embedded interest's GetDate
	return i.Interest.GetDate()
}

func (i *INRInterest) GetUSDAmount() float64 {
	return i.Amount
}

func (i *INRInterest) SetTTRate(rate float64) {
	i.TTRate = rate
}

func (i *INRInterest) SetTTDate(date time.Time) {
	i.TTDate = date
}

// Helper method for INR calculations
func (i *INRInterest) INRValue() float64 {
	return math.Round(i.Amount*i.TTRate*ROUNDING_FACTOR_2_DECIMALS) / ROUNDING_FACTOR_2_DECIMALS
}
