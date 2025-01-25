package fa

import "time"

// Broker statement transaction model
type Transaction struct {
	Security     string  `csv:"Security"`
	QuantitySold float64 `csv:"Quantity Sold"`
	DateAcquired string  `csv:"Date Acquired"`
	BuyingPrice  float64 `csv:"Buying Price (USD)"`
	DateSold     string  `csv:"Date Sold"`
	SellingPrice float64 `csv:"Selling Price (USD)"`
	Proceeds     float64 `csv:"Proceeds (USD)"`
	CostBasis    float64 `csv:"Cost Basis (USD)"`
	GainsLosses  float64 `csv:"Gains/Losses (USD)"`
}

// Position represents a snapshot of holdings at a point in time
type Position struct {
	Date     time.Time
	Quantity float64
	USDPrice float64
	USDValue float64
}

// PositionAnalysis tracks key positions for a ticker
type PositionAnalysis struct {
	Ticker          string
	FirstPosition   Position
	PeakPosition    Position
	YearEndPosition Position
}
