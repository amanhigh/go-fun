package tax

// DividendBase represents raw CSV data from broker statement
type DividendBase struct {
	Security         string  `csv:"Security"`
	DividendDate     string  `csv:"Dividend Date"`
	DividendPerShare float64 `csv:"Dividend Per Share (USD)"`
	DividendReceived float64 `csv:"Dividend Received (USD)"`
	DividendTax      float64 `csv:"Dividend Tax (USD)"`
	NetDividend      float64 `csv:"Net Dividend (USD)"`
}

// Dividend represents processed dividend data with INR conversions
type Dividend struct {
	DividendBase           // Embed input fields
	USDINRRate     float64 // TT Buy rate for conversion
	NetDividendINR float64 // Net amount in INR
	DividendTaxINR float64 // Tax amount in INR
}
