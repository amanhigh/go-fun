package tax

// StockData represents daily closing prices for a stock
// Map of date (YYYY-MM-DD format) to closing price
// This is the minimal format required for all tax calculations
type StockData struct {
	Prices map[string]float64 `json:"prices"` // Date (YYYY-MM-DD) -> Closing Price
}
