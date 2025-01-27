package tax

type PeakPrice struct {
	Ticker string `json:"ticker"`
	// Peak price information
	Date  string  `json:"peak_date"`
	Price float64 `json:"peak_price"`
}
