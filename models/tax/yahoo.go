package tax

// YahooChartResponse represents the complete structure of Yahoo Finance API response
type YahooChartResponse struct {
	Chart YahooChart `json:"chart"`
}

type YahooChart struct {
	Result []YahooChartResult `json:"result"`
	Error  any                `json:"error"`
}

// YahooChartResult contains all stock data for a ticker from Yahoo Finance
type YahooChartResult struct {
	Meta       YahooMeta       `json:"meta"`
	Timestamp  []int64         `json:"timestamp"`
	Indicators YahooIndicators `json:"indicators"`
	AdjClose   []float64       `json:"adjclose,omitempty"`
}

// YahooMeta contains metadata about the stock
type YahooMeta struct {
	Currency             string             `json:"currency"`
	Symbol               string             `json:"symbol"`
	ExchangeName         string             `json:"exchangeName"`
	FullExchangeName     string             `json:"fullExchangeName"`
	InstrumentType       string             `json:"instrumentType"`
	FirstTradeDate       int64              `json:"firstTradeDate"`
	RegularMarketTime    int64              `json:"regularMarketTime"`
	HasPrePostMarketData bool               `json:"hasPrePostMarketData"`
	GMTOffset            int                `json:"gmtoffset"`
	TimeZone             string             `json:"timezone"`
	ExchangeTimezoneName string             `json:"exchangeTimezoneName"`
	RegularMarketPrice   float64            `json:"regularMarketPrice"`
	FiftyTwoWeekHigh     float64            `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow      float64            `json:"fiftyTwoWeekLow"`
	RegularMarketDayHigh float64            `json:"regularMarketDayHigh"`
	RegularMarketDayLow  float64            `json:"regularMarketDayLow"`
	RegularMarketVolume  int64              `json:"regularMarketVolume"`
	LongName             string             `json:"longName"`
	ShortName            string             `json:"shortName"`
	ChartPreviousClose   float64            `json:"chartPreviousClose"`
	PriceHint            int                `json:"priceHint"`
	CurrentTradingPeriod YahooTradingPeriod `json:"currentTradingPeriod"`
	DataGranularity      string             `json:"dataGranularity"`
	Range                string             `json:"range"`
	ValidRanges          []string           `json:"validRanges"`
	// Legacy field name compatibility
	LastRefresh string `json:"lastRefresh,omitempty"`
}

// YahooTradingPeriod represents trading hours for a day
type YahooTradingPeriod struct {
	Pre     YahooPeriod `json:"pre"`
	Regular YahooPeriod `json:"regular"`
	Post    YahooPeriod `json:"post"`
}

// YahooPeriod represents a specific trading period
type YahooPeriod struct {
	Timezone  string `json:"timezone"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	GMTOffset int    `json:"gmtoffset"`
}

// YahooIndicators contains the actual price data
type YahooIndicators struct {
	Quote []YahooQuote `json:"quote"`
}

// YahooQuote contains OHLCV data for a single day
type YahooQuote struct {
	Open     []float64 `json:"open"`
	High     []float64 `json:"high"`
	Low      []float64 `json:"low"`
	Close    []float64 `json:"close"`
	Volume   []int64   `json:"volume"`
	AdjClose []float64 `json:"adjclose,omitempty"`
}
