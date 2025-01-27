package tax

// AlphaVantage response models
type VantageStockData struct {
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
