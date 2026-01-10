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
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/rs/zerolog/log"
)

// TickerManager manages ticker data lifecycle - download, cache, and price lookup.
type TickerManager interface {
	// DownloadTicker fetches ticker data from YahooClient and saves to file.
	DownloadTicker(ctx context.Context, ticker string) (err common.HttpError)

	// FindPeakPrice finds highest close price for a year from pre-downloaded ticker data.
	// Requires file to exist - call DownloadTicker(ctx, ticker) first if needed.
	// No caching or auto-download (unlike GetPrice which has both).
	FindPeakPrice(ctx context.Context, ticker string, year int) (tax.PeakPrice, common.HttpError)

	// GetPrice returns closing price for a date. Uses in-memory cache and auto-downloads if file missing.
	// Best for repeated calls on same ticker (valuation calculations).
	GetPrice(ctx context.Context, ticker string, date time.Time) (float64, common.HttpError)
}

type TickerManagerImpl struct {
	client    clients.StockDataClient
	downloads string
	cache     map[string]tax.StockData
	cacheLock sync.RWMutex
}

func NewTickerManager(client clients.StockDataClient, downloads string) *TickerManagerImpl {
	return &TickerManagerImpl{
		client:    client,
		downloads: downloads,
		cache:     make(map[string]tax.StockData),
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
	var data tax.StockData
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

func (t *TickerManagerImpl) FindPeakPrice(_ context.Context, ticker string, year int) (peakPrice tax.PeakPrice, err common.HttpError) {
	stockData, err := t.readTickerData(ticker)
	if err != nil {
		return peakPrice, err
	}

	yearStr := strconv.Itoa(year)
	return t.analyzePrices(stockData.Prices, ticker, yearStr), nil
}

func (t *TickerManagerImpl) GetPrice(ctx context.Context, ticker string, date time.Time) (float64, common.HttpError) {
	// Get cached/loaded data
	data, err := t.getTickerData(ctx, ticker)
	if err != nil {
		return 0, err
	}

	// Format date for lookup
	dateStr := date.Format(time.DateOnly)

	// Try exact date match first
	if closePrice, exists := data.Prices[dateStr]; exists {
		return closePrice, nil
	}

	// Find closest previous date if exact not found
	price, err := t.findClosestPreviousPrice(ticker, data, dateStr)
	if err == nil {
		return price, nil
	}

	return 0, common.NewHttpError("No price data found", http.StatusNotFound)
}

func (t *TickerManagerImpl) findClosestPreviousPrice(ticker string, data tax.StockData, dateStr string) (float64, common.HttpError) {
	var closestDate string
	for priceDate := range data.Prices {
		if priceDate <= dateStr && (closestDate == "" || priceDate > closestDate) {
			closestDate = priceDate
		}
	}

	if closestDate != "" {
		if closePrice, exists := data.Prices[closestDate]; exists {
			log.Debug().
				Str("Ticker", ticker).
				Str("RequestedDate", dateStr).
				Str("ClosestDate", closestDate).
				Float64("Price", closePrice).
				Msg("Using closest previous date price")
			return closePrice, nil
		}
	}
	return 0, common.NewHttpError("No price data found", http.StatusNotFound)
}

func (t *TickerManagerImpl) getTickerData(ctx context.Context, ticker string) (data tax.StockData, err common.HttpError) {
	// Try cache first
	t.cacheLock.RLock()
	data, exists := t.cache[ticker]
	t.cacheLock.RUnlock()

	if exists {
		log.Debug().Str("Ticker", ticker).Msg("Cache Hit")
		return data, nil
	}

	// Cache miss - load from file or auto-download
	data, err = t.loadOrDownloadTickerData(ctx, ticker)
	if err == nil {
		t.cacheLock.Lock()
		t.cache[ticker] = data
		t.cacheLock.Unlock()
		log.Debug().Str("Ticker", ticker).Msg("Added to Cache")
	}

	return
}

// loadOrDownloadTickerData attempts to load ticker data from file,
// and auto-downloads if file is missing
func (t *TickerManagerImpl) loadOrDownloadTickerData(ctx context.Context, ticker string) (tax.StockData, common.HttpError) {
	// Try loading from file first
	data, err := t.readTickerData(ticker)
	if err == nil {
		return data, nil
	}

	// File not found - attempt auto-download
	log.Info().Str("Ticker", ticker).Msg("Ticker file missing, attempting auto-download")

	if downloadErr := t.DownloadTicker(ctx, ticker); downloadErr != nil {
		log.Error().Str("Ticker", ticker).Err(downloadErr).Msg("Failed to auto-download ticker")
		return data, common.NewServerError(fmt.Errorf("failed to auto-download ticker %s: %w", ticker, downloadErr))
	}

	// Retry reading after download
	data, err = t.readTickerData(ticker)
	if err != nil {
		return data, common.NewServerError(fmt.Errorf("failed to read ticker %s after auto-download: %w", ticker, err))
	}

	log.Info().Str("Ticker", ticker).Msg("Successfully auto-downloaded and loaded ticker data")
	return data, nil
}

func (t *TickerManagerImpl) readTickerData(ticker string) (tax.StockData, common.HttpError) {
	var stockData tax.StockData

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

func (t *TickerManagerImpl) analyzePrices(prices map[string]float64, ticker, yearStr string) tax.PeakPrice {
	var highestClose float64
	var highestDate string

	for date, closePrice := range prices {
		if !strings.HasPrefix(date, yearStr) {
			continue
		}
		// Track highest close
		if closePrice > highestClose {
			highestClose = closePrice
			highestDate = date
		}
	}

	return tax.PeakPrice{
		Ticker: ticker,
		Date:   highestDate,
		Price:  highestClose,
	}
}
