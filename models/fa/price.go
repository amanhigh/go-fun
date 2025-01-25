package fa

// AlphaVantage response models
type StockData struct {
	MetaData   MetaData            `json:"Meta Data"`
	TimeSeries map[string]DayPrice `json:"Time Series (Daily)"`
}

type MetaData struct {
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	TimeZone      string `json:"5. Time Zone"`
}

type DayPrice struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

// SBI response models
type ExchangeRates struct {
	Rates []Rate
}

type Rate struct {
	Date   string  `json:"DATE"`
	TTBuy  float64 `json:"TT BUY"`
	TTSell float64 `json:"TT SELL"`
}

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

// File name constant for SBI Rate CSV
const SBI_RATES_FILENAME = "SBI_REFERENCE_RATES_USD.csv"
