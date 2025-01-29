package repository

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
)

type BaseCSVRepository[T any] interface {
	GetAllRecords(ctx context.Context) ([]T, common.HttpError)
	GetUniqueTickers(ctx context.Context) ([]string, common.HttpError)
	GetRecordsForTicker(ctx context.Context, ticker string) ([]T, common.HttpError)
}

type BaseCSVRepositoryImpl[T any] struct {
	filePath string
	records  []T          // Cache for records
	cache    sync.RWMutex // Lock for thread-safe access
}

func NewBaseCSVRepository[T any](filePath string) *BaseCSVRepositoryImpl[T] {
	return &BaseCSVRepositoryImpl[T]{
		filePath: filePath,
		records:  []T{},
	}
}

func (b *BaseCSVRepositoryImpl[T]) GetAllRecords(ctx context.Context) (records []T, err common.HttpError) {
	// Load records if needed
	if len(b.records) == 0 {
		if err = b.loadRecords(ctx); err != nil {
			return nil, err
		}
	}

	// Read Lock for accessing cache
	b.cache.RLock()
	records = b.records
	b.cache.RUnlock()

	return records, nil
}
func (b *BaseCSVRepositoryImpl[T]) GetUniqueTickers(ctx context.Context) (tickers []string, err common.HttpError) {
	// Get all records
	records, err := b.GetAllRecords(ctx)
	if err != nil {
		return nil, err
	}

	// Use map to track unique tickers
	tickerMap := make(map[string]bool)
	for _, record := range records {
		// Use type assertion to get Symbol field
		if ticker, ok := any(record).(tax.CSVRecord); ok {
			tickerMap[ticker.GetSymbol()] = true
		}
	}

	// Convert map keys to slice
	for ticker := range tickerMap {
		tickers = append(tickers, ticker)
	}

	return tickers, nil
}

func (b *BaseCSVRepositoryImpl[T]) GetRecordsForTicker(ctx context.Context, ticker string) (filtered []T, err common.HttpError) {
	// Get all records
	records, err := b.GetAllRecords(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by ticker
	for _, record := range records {
		if t, ok := any(record).(tax.CSVRecord); ok && t.GetSymbol() == ticker {
			filtered = append(filtered, record)
		}
	}

	return filtered, nil
}

func (b *BaseCSVRepositoryImpl[T]) loadRecords(ctx context.Context) (err common.HttpError) {
	// Lock for loading
	b.cache.Lock()
	defer b.cache.Unlock()

	// Double check after lock
	if len(b.records) > 0 {
		return nil
	}

	records, err := b.readCSVFile(ctx)
	if err != nil {
		return err
	}
	b.records = records
	return nil
}

func (b *BaseCSVRepositoryImpl[T]) readCSVFile(ctx context.Context) (records []T, err common.HttpError) {
	file, openErr := os.Open(b.filePath)
	if openErr != nil {
		log.Ctx(ctx).Error().Err(openErr).Str("path", b.filePath).Msg("Failed to open CSV file")
		return nil, common.NewServerError(openErr)
	}
	defer file.Close()

	if parseErr := gocsv.UnmarshalFile(file, &records); parseErr != nil {
		log.Ctx(ctx).Error().Err(parseErr).Msg("Failed to parse CSV")
		return nil, common.NewServerError(parseErr)
	}

	if len(records) == 0 {
		return nil, common.NewHttpError("empty CSV file", http.StatusBadRequest)
	}

	// Update type assertion to use CSVRecord
	if record, ok := any(records[0]).(tax.CSVRecord); ok {
		if !record.IsValid() {
			return nil, common.NewHttpError("invalid CSV format", http.StatusBadRequest)
		}
	}

	return records, nil
}
