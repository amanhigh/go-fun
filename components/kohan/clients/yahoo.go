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
	client *resty.Client
}

// NewYahooClient creates a new Yahoo Finance client
func NewYahooClient(client *resty.Client) *YahooClient {
	return &YahooClient{client: client}
}

// FetchDailyPrices fetches daily closing prices from Yahoo Finance
func (y *YahooClient) FetchDailyPrices(ctx context.Context, ticker string) (tax.StockData, common.HttpError) {
	var response tax.YahooChartResponse

	_, resErr := y.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"interval": "1d",
			"range":    "max",
		}).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		SetResult(&response).
		Get(fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", ticker))

	if err := util.ResponseProcessor(nil, resErr); err != nil {
		return tax.StockData{}, err
	}

	// Validate response structure
	if len(response.Chart.Result) == 0 {
		return tax.StockData{}, common.NewServerError(fmt.Errorf("no data found for %s", ticker))
	}

	result := response.Chart.Result[0]
	if len(result.Indicators.Quote) == 0 {
		return tax.StockData{}, common.NewServerError(fmt.Errorf("no quote data found for %s", ticker))
	}

	quote := result.Indicators.Quote[0]
	if len(quote.Close) == 0 {
		return tax.StockData{}, common.NewServerError(fmt.Errorf("no price data found for %s", ticker))
	}

	// Build prices map from timestamps and closing prices
	prices := make(map[string]float64)
	for i, ts := range result.Timestamp {
		date := time.Unix(ts, 0).UTC().Format(time.DateOnly)
		if i < len(quote.Close) {
			prices[date] = quote.Close[i]
		}
	}

	log.Info().
		Str("Ticker", ticker).
		Int("DataPoints", len(prices)).
		Msg("Successfully fetched ticker data from Yahoo Finance")

	return tax.StockData{Prices: prices}, nil
}
