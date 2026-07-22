package clients

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	tax "github.com/amanhigh/go-fun/models/tax"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

// SecurityClient defines the interface for fetching security information and price data.
type SecurityClient interface {
	FetchDailyPrices(ctx context.Context, ticker string) (tax.StockData, common.HttpError)
	GetSecurityInfo(ctx context.Context, query string) ([]tax.SecurityInfo, common.HttpError)
}

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

// Compile-time assertion: YahooClient implements SecurityClient
var _ SecurityClient = (*YahooClient)(nil)

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
			"events":   "splits",
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

	splits := y.extractSplits(&response)

	log.Info().
		Str("Ticker", ticker).
		Int("DataPoints", len(prices)).
		Int("Splits", len(splits)).
		Msg("Successfully fetched ticker data from Yahoo Finance")

	return tax.StockData{Prices: prices, Splits: splits}, nil
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

// extractSplits parses Yahoo's nested events.splits map into a chronologically ordered slice.
// Returns a non-nil empty slice when no split events exist.
func (y *YahooClient) extractSplits(response *tax.YahooChartResponse) []tax.YahooSplit {
	if len(response.Chart.Result) == 0 {
		return []tax.YahooSplit{}
	}

	result := response.Chart.Result[0]
	splitsMap := result.Events.Splits
	if splitsMap == nil {
		return []tax.YahooSplit{}
	}

	splits := make([]tax.YahooSplit, 0, len(splitsMap))
	for _, split := range splitsMap {
		splits = append(splits, split)
	}

	// Sort chronologically by Date
	sort.SliceStable(splits, func(i, j int) bool {
		return splits[i].Date < splits[j].Date
	})

	return splits
}

// GetSecurityInfo searches for securities by query string using Yahoo Finance /v1/finance/search.
// Returns all matching equity and ETF candidates without selecting one, or an empty non-nil slice.
func (y *YahooClient) GetSecurityInfo(ctx context.Context, query string) ([]tax.SecurityInfo, common.HttpError) {
	var response tax.YahooSearchResponse

	resp, resErr := y.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"q":           query,
			"quotesCount": "20",
			"newsCount":   "0",
		}).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		SetResult(&response).
		Get(y.getSearchURL())

	if err := util.ResponseProcessor(resp, resErr); err != nil {
		return nil, err
	}

	return y.extractSecurityInfo(&response), nil
}

// extractSecurityInfo parses Yahoo search quotes into SecurityInfo records.
// Filters to equity and ETF quote types, uses long name with short-name fallback.
// Returns an empty non-nil slice when there are no matching candidates.
func (y *YahooClient) extractSecurityInfo(response *tax.YahooSearchResponse) []tax.SecurityInfo {
	if len(response.Quotes) == 0 {
		return []tax.SecurityInfo{}
	}

	results := make([]tax.SecurityInfo, 0, len(response.Quotes))
	for _, quote := range response.Quotes {
		if quote.QuoteType != "EQUITY" && quote.QuoteType != "ETF" {
			continue
		}
		name := quote.LongName
		if name == "" {
			name = quote.ShortName
		}
		results = append(results, tax.SecurityInfo{
			Symbol:   quote.Symbol,
			Name:     name,
			Exchange: quote.Exchange,
			Type:     quote.QuoteType,
		})
	}

	if len(results) == 0 {
		return []tax.SecurityInfo{}
	}
	return results
}

// getSearchURL builds the search URL using the configured base URL
func (y *YahooClient) getSearchURL() string {
	return fmt.Sprintf("%s/v1/finance/search", y.baseURL)
}

// getChartURL builds the chart URL using the configured base URL
func (y *YahooClient) getChartURL(ticker string) string {
	return fmt.Sprintf("%s/v8/finance/chart/%s", y.baseURL, ticker)
}
