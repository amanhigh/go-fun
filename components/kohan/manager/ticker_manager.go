package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
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

	// GetPrice returns closing price for a date. Uses in-memory cache and auto-downloads if file missing.
	// Best for repeated calls on same ticker (valuation calculations).
	GetPrice(ctx context.Context, ticker string, date time.Time) (float64, common.HttpError)

	// GetDailyPrices returns all available closing prices for a given year as map[date]price.
	// Used for daily peak evaluation in valuation calculations.
	// Date format in returned map: "YYYY-MM-DD"
	GetDailyPrices(ctx context.Context, ticker string, year int) (map[string]float64, common.HttpError)

	// GetSplits returns split events within the given date range (inclusive).
	// Returns chronologically ordered defensive copy and non-nil empty slice.
	GetSplits(ctx context.Context, ticker string, from, to time.Time) ([]tax.SplitInfo, common.HttpError)
}

type TickerManagerImpl struct {
	client             clients.SecurityClient
	securityIDProvider SecurityIDProvider
	downloads          string
	cache              map[string]tax.StockData
	cacheLock          sync.RWMutex
}

func NewTickerManager(client clients.SecurityClient, securityIDProvider SecurityIDProvider, downloads string) *TickerManagerImpl {
	return &TickerManagerImpl{
		client:             client,
		securityIDProvider: securityIDProvider,
		downloads:          downloads,
		cache:              make(map[string]tax.StockData),
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

	// Fetch data using the private helper (tries direct fetch, then fallback resolution)
	var data tax.StockData
	if data, err = t.fetchTickerData(ctx, ticker); err != nil {
		return err
	}

	if err := t.saveTickerData(data, ticker); err != nil {
		return err
	}

	log.Info().Str("Ticker", ticker).Str("Path", filePath).Msg("Ticker data downloaded and saved")
	return nil
}

func (t *TickerManagerImpl) GetPrice(ctx context.Context, ticker string, date time.Time) (float64, common.HttpError) {
	data, err := t.getTickerData(ctx, ticker)
	if err != nil {
		return 0, err
	}

	dateStr := date.Format(time.DateOnly)

	var rawPrice float64
	if closePrice, exists := data.Prices[dateStr]; exists {
		rawPrice = closePrice
	} else if closestDate, closestPrice, ok := findClosestPriceOnOrBefore(data.Prices, dateStr); ok {
		log.Debug().
			Str("Ticker", ticker).
			Str("RequestedDate", dateStr).
			Str("ClosestDate", closestDate).
			Float64("Price", closestPrice).
			Msg("Using closest previous date price")
		rawPrice = closestPrice
	} else {
		return 0, common.NewHttpError("No price data found", http.StatusNotFound)
	}

	adjustedPrice, adjErr := t.adjustPriceForSplits(rawPrice, date, data.Splits, ticker)
	if adjErr != nil {
		return 0, adjErr
	}
	return adjustedPrice, nil
}

// adjustPriceForSplits returns the historical date's price expressed on the historical date's share basis
// by multiplying the raw cached price by cumulative split ratios
// for all split events on strictly later calendar dates.
func (t *TickerManagerImpl) adjustPriceForSplits(price float64, date time.Time, splits []tax.SplitInfo, ticker string) (float64, common.HttpError) {
	if vErr := validateSplits(splits, ticker); vErr != nil {
		return 0, vErr
	}
	return price * cumulativeSplitFactor(splits, date), nil
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
// auto-downloads if file is missing.
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

// saveTickerData marshals StockData and writes it to the ticker's JSON file.
func (t *TickerManagerImpl) saveTickerData(data tax.StockData, ticker string) common.HttpError {
	filePath := filepath.Join(t.downloads, ticker+".json")
	jsonData, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		return common.NewServerError(marshalErr)
	}
	if marshalErr = os.WriteFile(filePath, jsonData, util.DEFAULT_PERM); marshalErr != nil {
		return common.NewServerError(marshalErr)
	}
	return nil
}

// fetchTickerData tries the requested ticker first; on failure, resolves via
// SecurityIDProvider, returns common.ErrNotFound for zero candidates and
// common.ErrEntityExists for multiple candidates, and retries the sole candidate.
// Provider/search/retry errors are returned unchanged.
func (t *TickerManagerImpl) fetchTickerData(ctx context.Context, ticker string) (tax.StockData, common.HttpError) {
	// Try direct fetch first
	data, err := t.client.FetchDailyPrices(ctx, ticker)
	if err == nil {
		return data, nil
	}

	// Direct fetch failed — attempt fallback resolution via SecurityIDProvider
	securityID, secErr := t.securityIDProvider.GetSecurityID(ctx, ticker)
	if secErr != nil {
		return tax.StockData{}, secErr
	}

	candidates, searchErr := t.client.GetSecurityInfo(ctx, securityID)
	if searchErr != nil {
		return tax.StockData{}, searchErr
	}

	if len(candidates) == 0 {
		return tax.StockData{}, common.ErrNotFound
	}
	if len(candidates) > 1 {
		return tax.StockData{}, common.ErrEntityExists
	}

	// Exactly one candidate — retry with its symbol (unchanged, no trim/validation)
	return t.client.FetchDailyPrices(ctx, candidates[0].Symbol)
}

// GetSplits returns split events within the given date range (inclusive).
// Returns chronologically ordered defensive copy and non-nil empty slice.
func (t *TickerManagerImpl) GetSplits(ctx context.Context, ticker string, from, to time.Time) ([]tax.SplitInfo, common.HttpError) {
	if from.After(to) {
		return nil, common.NewHttpError("from date must be before or equal to to date", http.StatusBadRequest)
	}

	data, err := t.getTickerData(ctx, ticker)
	if err != nil {
		return nil, err
	}

	if vErr := validateSplits(data.Splits, ticker); vErr != nil {
		return nil, vErr
	}

	fromDay := from.UTC().Truncate(24 * time.Hour).Unix() //nolint:mnd
	toDay := to.UTC().Truncate(24 * time.Hour).Unix()     //nolint:mnd

	result := make([]tax.SplitInfo, 0)
	for _, split := range data.Splits {
		splitDay := split.EffectiveDate().Unix()
		if splitDay >= fromDay && splitDay <= toDay {
			result = append(result, split)
		}
	}

	// Enforce chronological order regardless of storage order
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	return result, nil
}

// filterAndAdjustYearPrices filters prices for a given year and applies split adjustments.
func (t *TickerManagerImpl) filterAndAdjustYearPrices(data tax.StockData, year int, ticker string) (map[string]float64, common.HttpError) {
	yearStr := strconv.Itoa(year)
	yearPrices := make(map[string]float64)
	for date, price := range data.Prices {
		if strings.HasPrefix(date, yearStr) {
			parsedDate, parseErr := time.Parse(time.DateOnly, date)
			if parseErr == nil {
				adjustedPrice, adjErr := t.adjustPriceForSplits(price, parsedDate, data.Splits, ticker)
				if adjErr != nil {
					return nil, adjErr
				}
				yearPrices[date] = adjustedPrice
			} else {
				yearPrices[date] = price
			}
		}
	}
	return yearPrices, nil
}

// findClosestPriceOnOrBefore returns the closest date on or before refDate, its price, and whether a match was found.
func findClosestPriceOnOrBefore(prices map[string]float64, refDate string) (closestDate string, price float64, ok bool) {
	for date := range prices {
		if date <= refDate && (closestDate == "" || date > closestDate) {
			closestDate = date
		}
	}
	if closestDate == "" {
		return "", 0, false
	}
	price, ok = prices[closestDate]
	return
}

// GetDailyPrices returns all available closing prices for a given year
// as a map with date format "YYYY-MM-DD" as key and price as value.
// Used for daily peak INR value evaluation during valuation calculations.
func (t *TickerManagerImpl) GetDailyPrices(ctx context.Context, ticker string, year int) (map[string]float64, common.HttpError) {
	data, err := t.getTickerData(ctx, ticker)
	if err != nil {
		return nil, err
	}

	yearPrices, filterErr := t.filterAndAdjustYearPrices(data, year, ticker)
	if filterErr != nil {
		return nil, filterErr
	}

	// Return error if no prices found for the requested year
	if len(yearPrices) == 0 {
		return nil, common.NewHttpError(
			fmt.Sprintf("no price data found for ticker %s in year %d", ticker, year),
			http.StatusNotFound,
		)
	}

	// missing exact split-date prices are unsafe to synthesize because
	// the post-split share basis is unverified; investigate and add a
	// verified fallback only if real data requires it.
	for _, split := range data.Splits {
		if splitDateStr := split.EffectiveDate().Format(time.DateOnly); split.EffectiveDate().Year() == year {
			if _, exists := data.Prices[splitDateStr]; !exists {
				return nil, common.NewHttpError(
					fmt.Sprintf("no cached price on split date %s for ticker %s", splitDateStr, ticker),
					http.StatusNotFound,
				)
			}
		}
	}

	// Backfill previous year-end price with split adjustment for carry-over.
	prevYearEnd := fmt.Sprintf("%d-12-31", year-1)
	if prevPrice, exists := data.Prices[prevYearEnd]; exists {
		prevDate := time.Date(year-1, 12, 31, 0, 0, 0, 0, time.UTC)
		adjustedPrice, adjErr := t.adjustPriceForSplits(prevPrice, prevDate, data.Splits, ticker)
		if adjErr != nil {
			return nil, adjErr
		}
		yearPrices[prevYearEnd] = adjustedPrice
	}

	return yearPrices, nil
}
