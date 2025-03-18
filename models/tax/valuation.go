package tax

import "time"

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

func (t Trade) GetKey() string {
	return t.Symbol
}

func (t Trade) GetDate() time.Time {
	date, _ := time.Parse(time.DateOnly, t.Date)
	return date
}

func (t Trade) IsValid() bool {
	return t.Symbol != "" && t.Date != "" && t.Type != "" &&
		(t.Type == "BUY" || t.Type == "SELL")
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

// INRValutaion mirrors Valuation structure with tax positions
type INRValutaion struct {
	Ticker          string
	FirstPosition   INRPosition // First position with exchange rate details
	PeakPosition    INRPosition // Peak position with exchange rate details
	YearEndPosition INRPosition // Year end position with exchange rate details
}

// Helper to create tax valuation from base valuation
func NewINRValuation(valuation Valuation) INRValutaion {
	return INRValutaion{
		Ticker:          valuation.Ticker,
		FirstPosition:   INRPosition{Position: valuation.FirstPosition},
		PeakPosition:    INRPosition{Position: valuation.PeakPosition},
		YearEndPosition: INRPosition{Position: valuation.YearEndPosition},
	}
}

// INRPosition extends Position with exchange rate details
// INRPosition is not a CSV Record but Implements Exchangeable Interface.
type INRPosition struct {
	Position           // Embed original position
	TTDate   time.Time // Date for which exchange rate is applied
	TTRate   float64   // TT Buy rate used for conversion
}

// INRValue calculates INR value using embedded position's USD value
// BUG: Should INRValue be part of Interface or remove if unused.
func (t *INRPosition) INRValue() float64 {
	return t.USDValue() * t.TTRate
}

// Implement Exchangeable interface for INRPosition
func (t *INRPosition) GetDate() time.Time {
	return t.Position.Date
}

func (t *INRPosition) GetKey() string {
	// This mostly won't be used is to Satisfy CSV Record.
	return t.Position.Date.Format(time.DateOnly)
}

func (t *INRPosition) GetUSDAmount() float64 {
	return t.USDValue()
}

func (t *INRPosition) SetTTRate(rate float64) {
	t.TTRate = rate
}

func (t *INRPosition) SetTTDate(date time.Time) {
	t.TTDate = date
}

func (t *INRPosition) IsValid() bool {
	return !t.Position.Date.IsZero() && t.TTRate > 0
}
