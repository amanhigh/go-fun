package tax

// Account represents year-end account position
type Account struct {
	Symbol      string  `csv:"Symbol"`
	Quantity    float64 `csv:"Quantity"`
	Cost        float64 `csv:"Cost"`
	MarketValue float64 `csv:"MarketValue"`
}

// GetSymbol implements CSVRecord interface
func (a Account) GetSymbol() string {
	return a.Symbol
}

// IsValid implements CSVRecord interface
func (a Account) IsValid() bool {
	return a.Symbol != "" && a.Quantity != 0 && a.Cost != 0 && a.MarketValue != 0
}