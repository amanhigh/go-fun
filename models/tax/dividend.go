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
