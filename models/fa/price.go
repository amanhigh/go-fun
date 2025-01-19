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
