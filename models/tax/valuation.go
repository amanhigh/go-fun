package tax

import (
	"fmt"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

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

func (t Trade) GetDate() (time.Time, common.HttpError) {
	parsedTime, err := time.Parse(time.DateOnly, t.Date)
	if err != nil {
		return time.Time{}, NewInvalidDateError(fmt.Sprintf("failed to parse date '%s': %v", t.Date, err))
	}
	return parsedTime, nil
}

func (t Trade) IsValid() bool {
	if t.Symbol == "" || t.Date == "" || t.Type == "" {
		return false
	}

	// Accept both uppercase and mixed case trade types (BUY/Buy, SELL/Sell)
	uppercaseType := strings.ToUpper(t.Type)
	return uppercaseType == "BUY" || uppercaseType == "SELL"
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

// INRValuation mirrors Valuation structure with tax positions
type INRValuation struct {
	Ticker          string
	FirstPosition   INRPosition // First position with exchange rate details
	PeakPosition    INRPosition // Peak position with exchange rate details
	YearEndPosition INRPosition // Year end position with exchange rate details
}

// Helper to create tax valuation from base valuation
func NewINRValuation(valuation Valuation) INRValuation {
	return INRValuation{
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

// Implement Exchangeable interface for INRPosition
func (t *INRPosition) GetDate() (time.Time, common.HttpError) {
	// Check if the embedded Position Date is zero
	if t.Date.IsZero() {
		// Return an error indicating an invalid date
		return time.Time{}, NewInvalidDateError("position date is zero")
	}
	return t.Date, nil
}

func (t *INRPosition) GetKey() string {
	// This mostly won't be used is to Satisfy CSV Record.
	return t.Date.Format(time.DateOnly)
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
	return !t.Date.IsZero() && t.TTRate > 0
}
func (t *INRPosition) INRValue() float64 {
	return t.USDValue() * t.TTRate
}
