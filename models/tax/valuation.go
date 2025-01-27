package tax

import "time"

// Replace Symbolic interface with CSVRecord
type CSVRecord interface {
	GetSymbol() string // For ticker related functions
	IsValid() bool     // For CSV validation
}

// Broker statement trade model
type Trade struct {
	Symbol     string  `csv:"Symbol"`
	Date       string  `csv:"Date"`
	Type       string  `csv:"Type"`
	Quantity   float64 `csv:"Quantity"`
	USDPrice   float64 `csv:"Price"`
	USDValue   float64 `csv:"Value"`
	Commission float64 `csv:"Commission"`
}

func NewTrade(symbol, date, tradeType string, quantity, price float64) Trade {
	return Trade{
		Symbol:   symbol,
		Date:     date,
		Type:     tradeType,
		Quantity: quantity,
		USDPrice: price,
	}
}

// Add GetSymbol method to Trade struct
func (t Trade) GetSymbol() string {
	return t.Symbol
}

// Add IsValid implementation
func (t Trade) IsValid() bool {
	return t.Symbol != "" && t.Date != "" && t.Type != "" &&
		(t.Type == "BUY" || t.Type == "SELL")
}

func (t Trade) ParseDate() (time.Time, error) {
	return time.Parse(time.DateOnly, t.Date)
}

// Position represents a snapshot of holdings at a point in time
type Position struct {
	Date     time.Time
	Quantity float64
	USDPrice float64
}

func (p *Position) USDValue() float64 {
	return p.Quantity * p.USDPrice
}

// Valuation tracks key positions for a ticker
type Valuation struct {
	Ticker          string
	FirstPosition   Position
	PeakPosition    Position
	YearEndPosition Position
}
