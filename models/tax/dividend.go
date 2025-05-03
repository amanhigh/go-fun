package tax

import (
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

type Dividend struct {
	Symbol string  `csv:"Symbol"`
	Date   string  `csv:"Date"`
	Amount float64 `csv:"Amount"`
	Tax    float64 `csv:"Tax"`
	Net    float64 `csv:"Net"`
}

func (d Dividend) GetKey() string {
	return d.Symbol
}

func (d Dividend) IsValid() bool {
	return d.Symbol != "" && d.Date != "" && d.Amount != 0
}

func (d Dividend) GetDate() (time.Time, common.HttpError) {
	t, err := time.Parse(time.DateOnly, d.Date)
	if err != nil {
		return time.Time{}, NewInvalidDateError(fmt.Sprintf("failed to parse date '%s': %v", d.Date, err))
	}
	return t, nil
}

// INRDividend adds exchange rate details to basic dividend
type INRDividend struct {
	Dividend           // Embed original dividend
	TTDate   time.Time // Date for exchange rate (keeping for interface consistency)
	TTRate   float64   // Exchange rate
}

// Implement Exchangeable interface
func (d *INRDividend) GetDate() (time.Time, common.HttpError) {
	// Use embedded dividend's GetDate
	return d.Dividend.GetDate()
}

func (d *INRDividend) GetUSDAmount() float64 {
	return d.Amount // Using gross amount for conversion
}

func (d *INRDividend) SetTTRate(rate float64) {
	d.TTRate = rate
}

func (d *INRDividend) SetTTDate(date time.Time) {
	d.TTDate = date
}

// Helper method for INR calculations
func (d *INRDividend) INRValue() float64 {
	return d.Amount * d.TTRate
}
