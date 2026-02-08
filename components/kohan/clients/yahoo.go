package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	tax "github.com/amanhigh/go-fun/models/tax"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

// YahooClient fetches stock prices from Yahoo Finance API
type YahooClient struct {
	client              *resty.Client
	baseURL             string
	tickerDataStartYear int
}

// NewYahooClient creates a new Yahoo Finance client with custom base URL and ticker data start year
func NewYahooClient(client *resty.Client, baseURL string, tickerDataStartYear int) *YahooClient {
	return &YahooClient{
		client:              client,
		baseURL:             baseURL,
		tickerDataStartYear: tickerDataStartYear,
	}
}

// FetchDailyPrices fetches daily closing prices from Yahoo Finance
// Uses period1 and period2 parameters instead of "range" to get daily granularity
// "range": "max" returns monthly data, while period parameters return true daily data
func (y *YahooClient) FetchDailyPrices(ctx context.Context, ticker string) (tax.StockData, common.HttpError) {
	var response tax.YahooChartResponse

	// Use epoch timestamps for maximum date range
	// period1: based on TickerDataStartYear configuration (avoids sparse/missing data)
	// period2: now (today's date)
	// Note: Using period1/period2 returns true daily data, while "range": "max" returns only monthly
	startDate := time.Date(y.tickerDataStartYear, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	now := time.Now().Unix()

	resp, resErr := y.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"interval": "1d",
			"period1":  fmt.Sprintf("%d", startDate),
			"period2":  fmt.Sprintf("%d", now),
		}).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		SetResult(&response).
		Get(y.getChartURL(ticker))

	if err := util.ResponseProcessor(resp, resErr); err != nil {
		return tax.StockData{}, err
	}

	prices, err := y.extractPrices(&response, ticker)
	if err != nil {
		return tax.StockData{}, err
	}

	log.Info().
		Str("Ticker", ticker).
		Int("DataPoints", len(prices)).
		Msg("Successfully fetched ticker data from Yahoo Finance")

	return tax.StockData{Prices: prices}, nil
}

// extractPrices validates response and builds prices map from timestamps and closing prices
func (y *YahooClient) extractPrices(response *tax.YahooChartResponse, ticker string) (map[string]float64, common.HttpError) {
	if len(response.Chart.Result) == 0 {
		return nil, common.NewServerError(fmt.Errorf("no data found for %s", ticker))
	}

	result := response.Chart.Result[0]
	if len(result.Indicators.Quote) == 0 {
		return nil, common.NewServerError(fmt.Errorf("no quote data found for %s", ticker))
	}

	quote := result.Indicators.Quote[0]
	if len(quote.Close) == 0 {
		return nil, common.NewServerError(fmt.Errorf("no price data found for %s", ticker))
	}

	prices := make(map[string]float64)
	for i, ts := range result.Timestamp {
		date := time.Unix(ts, 0).UTC().Format(time.DateOnly)
		if i < len(quote.Close) {
			prices[date] = quote.Close[i]
		}
	}

	return prices, nil
}

// getChartURL builds the chart URL using the configured base URL
func (y *YahooClient) getChartURL(ticker string) string {
	return fmt.Sprintf("%s/v8/finance/chart/%s", y.baseURL, ticker)
}
