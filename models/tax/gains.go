package tax

import "time"

type Gains struct {
	Symbol     string  `csv:"Symbol"`
	BuyDate    string  `csv:"BuyDate"`
	SellDate   string  `csv:"SellDate"`
	Quantity   float64 `csv:"Quantity"`
	PNL        float64 `csv:"PNL"`
	Commission float64 `csv:"Commission"`
	Type       string  `csv:"Type"`
}

func (g Gains) GetSymbol() string {
	return g.Symbol
}

func (g Gains) IsValid() bool {
	return g.Symbol != "" && g.BuyDate != "" && g.SellDate != ""
}

func (g Gains) ParseBuyDate() (time.Time, error) {
	return time.Parse(time.DateOnly, g.BuyDate)
}

func (g Gains) ParseSellDate() (time.Time, error) {
	return time.Parse(time.DateOnly, g.SellDate)
}
