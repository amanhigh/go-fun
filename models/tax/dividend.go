package tax

import "time"

type Dividend struct {
	Symbol string  `csv:"Symbol"`
	Date   string  `csv:"Date"`
	Amount float64 `csv:"Amount"`
	Tax    float64 `csv:"Tax"`
	Net    float64 `csv:"Net"`
}

func (d Dividend) GetSymbol() string {
	return d.Symbol
}

func (d Dividend) IsValid() bool {
	return d.Symbol != "" && d.Date != "" && d.Amount != 0
}

func (d Dividend) ParseDate() (time.Time, error) {
	return time.Parse(time.DateOnly, d.Date)
}

// INRDividend adds exchange rate details to basic dividend
type INRDividend struct {
	Dividend           // Embed original dividend
	TTDate   time.Time // Date for exchange rate (keeping for interface consistency)
	TTRate   float64   // Exchange rate
}

// Implement Exchangeable interface
func (d *INRDividend) GetDate() time.Time {
	// Use original dividend date
	date, _ := d.ParseDate()
	return date
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
