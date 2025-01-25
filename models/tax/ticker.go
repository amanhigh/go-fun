package tax

// TickerAnalysis represents analyzed ticker data for a given year
type TickerAnalysis struct {
	Ticker string `json:"ticker"`
	// Peak price information
	PeakDate  string  `json:"peak_date"`
	PeakPrice float64 `json:"peak_price"`

	// Year end price information
	YearEndDate  string  `json:"year_end_date"`
	YearEndPrice float64 `json:"year_end_price"`
}

// TickerInfo extends TickerAnalysis with TT rate conversions
type TickerInfo struct {
	TickerAnalysis          // Embed base USD analysis
	PeakTTRate      float64 `json:"peak_tt_rate"`
	YearEndTTRate   float64 `json:"year_end_tt_rate"`
	PeakPriceINR    float64 `json:"peak_price_inr"`
	YearEndPriceINR float64 `json:"year_end_price_inr"`
}
