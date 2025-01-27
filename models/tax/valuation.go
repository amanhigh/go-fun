package tax

import "time"

// Broker statement trade model
type Trade struct {
	Symbol     string    // Stock symbol (e.g. MPC)
	Date       time.Time // Trade date
	Type       string    // BUY/SELL
	Quantity   float64   // Number of shares
	USDPrice   float64   // Price per share in USD
	USDValue   float64   // Total value (Quantity * Price)
	Commission float64   // Trade commission
}

// Position represents a snapshot of holdings at a point in time
type Position struct {
	Date     time.Time
	Quantity float64
	USDPrice float64
	USDValue float64
}

// Valuation tracks key positions for a ticker
type Valuation struct {
	Ticker          string
	FirstPosition   Position
	PeakPosition    Position
	YearEndPosition Position
}
