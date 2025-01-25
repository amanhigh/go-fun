package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fa"
	"github.com/rs/zerolog/log"
)

//go:generate mockery --name TickerManager
type TickerManager interface {
	DownloadTicker(ctx context.Context, ticker string) (err common.HttpError)
	// BUG: Rename yearly analysis
	AnalyzeTicker(ctx context.Context, ticker string, year int) (analysis fa.TickerAnalysis, err common.HttpError)
	// BUG: Rename GetPrice
	GetPriceOnDate(ctx context.Context, ticker string, date time.Time) (float64, error)
}

type TickerManagerImpl struct {
	client    clients.AlphaClient
	downloads string
	cache     map[string]fa.StockData
	cacheLock sync.RWMutex
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
	var data fa.StockData
	if data, err = t.client.FetchDailyPrices(ctx, ticker); err != nil {
		return err
	}

	// Save data to file
	if jsonData, marshalErr := json.Marshal(data); marshalErr == nil {
		if marshalErr = os.WriteFile(filePath, jsonData, util.DEFAULT_PERM); marshalErr != nil {
			return common.NewServerError(marshalErr)
		}
		log.Info().Str("Ticker", ticker).Str("Path", filePath).Msg("Ticker data downloaded and saved")
	} else {
		return common.NewServerError(marshalErr)
	}

	return nil
}

func (t *TickerManagerImpl) AnalyzeTicker(ctx context.Context, ticker string, year int) (analysis fa.TickerAnalysis, err common.HttpError) {
	stockData, err := t.readTickerData(ticker)
	if err != nil {
		return analysis, err
	}

	yearStr := strconv.Itoa(year)
	yearEndDate := fmt.Sprintf("%s-12-31", yearStr)

	return t.analyzeTimeSeries(stockData.TimeSeries, ticker, yearStr, yearEndDate), nil
}

func (t *TickerManagerImpl) GetPriceOnDate(ctx context.Context, ticker string, date time.Time) (float64, error) {
	// Get cached/loaded data
	data, err := t.getTickerData(ctx, ticker)
	if err != nil {
		return 0, err
	}

	// Format date for lookup
	// BUG: Use Date Constant in Models
	dateStr := date.Format("2006-01-02")

	// Try exact date match first
	if dayPrice, exists := data.TimeSeries[dateStr]; exists {
		if price, err := strconv.ParseFloat(dayPrice.Close, 64); err == nil {
			return price, nil
		}
	}

	// Find closest previous date if exact not found
	var closestDate string
	for tsDate := range data.TimeSeries {
		if tsDate <= dateStr && (closestDate == "" || tsDate > closestDate) {
			closestDate = tsDate
		}
	}

	if closestDate != "" {
		if dayPrice, exists := data.TimeSeries[closestDate]; exists {
			if price, err := strconv.ParseFloat(dayPrice.Close, 64); err == nil {
				log.Debug().
					Str("Ticker", ticker).
					Str("RequestedDate", dateStr).
					Str("ClosestDate", closestDate).
					Float64("Price", price).
					Msg("Using closest previous date price")
				return price, nil
			}
		}
	}

	return 0, common.NewHttpError("No price data found", http.StatusNotFound)
}

func (t *TickerManagerImpl) getTickerData(ctx context.Context, ticker string) (data fa.StockData, err common.HttpError) {
	// Try cache first
	t.cacheLock.RLock()
	data, exists := t.cache[ticker]
	t.cacheLock.RUnlock()

	if exists {
		log.Debug().Str("Ticker", ticker).Msg("Cache Hit")
		return data, nil
	}

	// Cache miss - load from file
	data, err = t.readTickerData(ticker)
	if err == nil {
		t.cacheLock.Lock()
		t.cache[ticker] = data
		t.cacheLock.Unlock()
		log.Debug().Str("Ticker", ticker).Msg("Added to Cache")
	}

	return
}

func (t *TickerManagerImpl) readTickerData(ticker string) (fa.StockData, common.HttpError) {
	var stockData fa.StockData

	filePath := filepath.Join(t.downloads, fmt.Sprintf("%s.json", ticker))
	data, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return stockData, common.NewServerError(readErr)
	}

	parseErr := json.Unmarshal(data, &stockData)
	if parseErr != nil {
		return stockData, common.NewServerError(parseErr)
	}

	return stockData, nil
}

func (t *TickerManagerImpl) analyzeTimeSeries(timeSeries map[string]fa.DayPrice, ticker, yearStr, yearEndDate string) fa.TickerAnalysis {
	var highestClose float64
	var highestDate string
	var yearEndClose float64
	var lastTradingDay string

	for date, values := range timeSeries {
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

	return fa.TickerAnalysis{
		Ticker:       ticker,
		PeakDate:     highestDate,
		PeakPrice:    highestClose,
		YearEndDate:  lastTradingDay,
		YearEndPrice: yearEndClose,
	}
}
