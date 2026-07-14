package tax

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

// YahooChartResponse represents the complete structure of Yahoo Finance API response
type YahooChartResponse struct {
	Chart YahooChart `json:"chart"`
}

type YahooChart struct {
	Result []YahooChartResult `json:"result"`
	Error  any                `json:"error"`
}

// YahooSplit represents a stock split event from the Yahoo Finance API
type YahooSplit struct {
	Date        int64   `json:"date"`
	Numerator   float64 `json:"numerator"`
	Denominator float64 `json:"denominator"`
}

// Validate checks that a YahooSplit has valid event data:
// positive finite numerator and denominator, and a valid usable Unix timestamp.
// The error message includes the ticker and event timestamp for traceability.
func (s YahooSplit) Validate(ticker string) common.HttpError {
	prefix := fmt.Sprintf("ticker %s: split event timestamp %d", ticker, s.Date)
	if s.Date <= 0 {
		return common.NewHttpError(fmt.Sprintf("%s: non-positive timestamp %d", prefix, s.Date), http.StatusBadRequest)
	}
	if s.Numerator <= 0 || math.IsInf(s.Numerator, 0) || math.IsNaN(s.Numerator) {
		return common.NewHttpError(fmt.Sprintf("%s: non-positive or non-finite numerator %f", prefix, s.Numerator), http.StatusBadRequest)
	}
	if s.Denominator <= 0 || math.IsInf(s.Denominator, 0) || math.IsNaN(s.Denominator) {
		return common.NewHttpError(fmt.Sprintf("%s: non-positive or non-finite denominator %f", prefix, s.Denominator), http.StatusBadRequest)
	}
	return nil
}

// EffectiveDate returns the UTC calendar date of the split at midnight,
// normalizing the Yahoo Unix timestamp to the start of its UTC day.
func (s YahooSplit) EffectiveDate() time.Time {
	return time.Unix(s.Date, 0).UTC().Truncate(24 * time.Hour) //nolint:mnd
}

// YahooEvents contains optional event data from the Yahoo Finance API
type YahooEvents struct {
	Splits map[string]YahooSplit `json:"splits"`
}

// YahooChartResult contains all stock data for a ticker from Yahoo Finance
type YahooChartResult struct {
	Meta       YahooMeta       `json:"meta"`
	Timestamp  []int64         `json:"timestamp"`
	Indicators YahooIndicators `json:"indicators"`
	AdjClose   []float64       `json:"adjclose,omitempty"`
	Events     YahooEvents     `json:"events"`
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
