package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fa"
	"github.com/rs/zerolog/log"
)

type TickerManager interface {
	DownloadTicker(ctx context.Context, ticker string) (err common.HttpError)
	AnalyzeTicker(ctx context.Context, ticker string, year int) (analysis fa.TickerAnalysis, err common.HttpError)
}

type TickerManagerImpl struct {
	client    clients.AlphaClient
	downloads string
}

func NewTickerManager(client clients.AlphaClient, downloads string) *TickerManagerImpl {
	return &TickerManagerImpl{
		client:    client,
		downloads: downloads,
	}
}

func (t *TickerManagerImpl) DownloadTicker(ctx context.Context, ticker string) (err common.HttpError) {
	filePath := filepath.Join(t.downloads, ticker+".json")

	// Check if file exists
	if info, statErr := os.Stat(filePath); statErr == nil {
		modTime := info.ModTime().Format("2006-01-02")
		log.Info().Str("Ticker", ticker).Str("ModTime", modTime).Msg("Ticker data already exists")
		return nil
	}

	// Create directory if it doesn't exist
	if err1 := os.MkdirAll(t.downloads, util.DIR_DEFAULT_PERM); err1 != nil {
		return common.NewServerError(err1)
	}

	// Fetch data using AlphaClient
	var data interface{}
	if data, err = t.client.FetchDailyPrices(ctx, ticker); err != nil {
		return err
	}

	// Save data to file
	if jsonData, err1 := json.Marshal(data); err1 == nil {
		if err1 = os.WriteFile(filePath, jsonData, util.DEFAULT_PERM); err1 != nil {
			return common.NewServerError(err1)
		}
		log.Info().Str("Ticker", ticker).Str("Path", filePath).Msg("Ticker data downloaded and saved")
	} else {
		return common.NewServerError(err1)
	}

	return nil
}

func (t *TickerManagerImpl) AnalyzeTicker(ctx context.Context, ticker string, year int) (analysis fa.TickerAnalysis, err common.HttpError) {
	// Read ticker file
	filePath := filepath.Join(t.downloads, fmt.Sprintf("%s.json", ticker))
	data, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return analysis, common.NewServerError(readErr)
	}

	// Parse stock data
	var stockData fa.StockData
	parseErr := json.Unmarshal(data, &stockData)
	if parseErr != nil {
		return analysis, common.NewServerError(parseErr)
	}

	// Initialize variables for analysis
	yearStr := strconv.Itoa(year)
	yearEndDate := fmt.Sprintf("%s-12-31", yearStr)

	var highestClose float64
	var highestDate string
	var yearEndClose float64
	var lastTradingDay string

	// Analyze time series data
	for date, values := range stockData.TimeSeries {
		// Only process dates for requested year
		if !strings.HasPrefix(date, yearStr) {
			continue
		}

		// Parse close price
		closePrice, parseErr := strconv.ParseFloat(values.Close, 64)
		if parseErr != nil {
			continue
		}

		// Track highest close
		if closePrice > highestClose {
			highestClose = closePrice
			highestDate = date
		}

		// Track year end close
		if date == yearEndDate {
			yearEndClose = closePrice
			lastTradingDay = date
		}

		// Keep track of last trading day
		if lastTradingDay == "" || date > lastTradingDay {
			lastTradingDay = date
			yearEndClose = closePrice
		}
	}

	// Prepare analysis result
	analysis = fa.TickerAnalysis{
		PeakDate:     highestDate,
		PeakPrice:    highestClose,
		YearEndDate:  lastTradingDay,
		YearEndPrice: yearEndClose,
	}

	return analysis, nil
}
