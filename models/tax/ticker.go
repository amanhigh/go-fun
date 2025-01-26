package tax

// BaseValuation represents analyzed ticker data for a given year
type BaseValuation struct {
	Ticker string `json:"ticker"`
	// Peak price information
	PeakDate  string  `json:"peak_date"`
	PeakPrice float64 `json:"peak_price"`

	// Year end price information
	YearEndDate  string  `json:"year_end_date"`
	YearEndPrice float64 `json:"year_end_price"`
}
