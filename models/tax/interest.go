package tax

import "time"

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

func (i Interest) GetDate() time.Time {
	date, _ := time.Parse(time.DateOnly, i.Date)
	return date
}

// INRInterest adds exchange rate details to basic interest
type INRInterest struct {
	Interest           // Embed original interest
	TTDate   time.Time // Date for exchange rate
	TTRate   float64   // Exchange rate
}

// Implement Exchangeable interface
func (i *INRInterest) GetDate() time.Time {
	date, _ := time.Parse(time.DateOnly, i.Date)
	return date
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
	return i.Amount * i.TTRate
}
