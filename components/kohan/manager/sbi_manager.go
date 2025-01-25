package manager

import (
	"context"
	"encoding/csv"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

//go:generate mockery --name SBIManager
type SBIManager interface {
	DownloadRates(ctx context.Context) common.HttpError
	GetTTBuyRate(date time.Time) (float64, common.HttpError)
}

type SBIManagerImpl struct {
	client      clients.SBIClient
	downloadDir string
	rateCache   map[string]float64
}

func NewSBIManager(client clients.SBIClient, downloadDir string) *SBIManagerImpl {
	return &SBIManagerImpl{
		client:      client,
		downloadDir: downloadDir,
		rateCache:   make(map[string]float64),
	}
}

// SaveRates fetches the latest exchange rates and saves them to a CSV file
// in the configured download directory. Returns an error if any step fails.
func (s *SBIManagerImpl) GetTTBuyRate(date time.Time) (rate float64, err common.HttpError) {
	if err = s.loadRatesIfNeeded(); err != nil {
		return 0, err
	}

	// BUG: Declare constant for date format in models
	dateStr := date.Format("2006-01-02")
	if rate, exists := s.rateCache[dateStr]; exists {
		return rate, nil
	}
	return 0, common.ErrNotFound
}

func (s *SBIManagerImpl) DownloadRates(ctx context.Context) (err common.HttpError) {
	// BUG: #B Skip dowloading if file already exists
	var csvContent string
	if csvContent, err = s.client.FetchExchangeRates(ctx); err != nil {
		return
	}

	// Ensure directory exists
	if err1 := os.MkdirAll(s.downloadDir, os.ModePerm); err1 != nil {
		return common.NewServerError(err1)
	}

	// Write to file
	filePath := filepath.Join(s.downloadDir, tax.SBI_RATES_FILENAME)
	if err1 := os.WriteFile(filePath, []byte(csvContent), util.DEFAULT_PERM); err1 != nil {
		return common.NewServerError(err1)
	}

	return
}

func (s *SBIManagerImpl) parseCSVToRateMap(records [][]string) (map[string]float64, common.HttpError) {
	rateMap := make(map[string]float64)
	for _, record := range records {
		// Parse date and rate
		dateStr := strings.Split(record[0], " ")[0] // Get date part only
		rate, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, common.NewServerError(err)
		}
		rateMap[dateStr] = rate
	}
	return rateMap, nil
}

func (s *SBIManagerImpl) loadRatesIfNeeded() common.HttpError {
	if len(s.rateCache) > 0 {
		return nil
	}

	// BUG: Private method to check file presence
	filePath := filepath.Join(s.downloadDir, tax.SBI_RATES_FILENAME)
	if _, err := os.Stat(filePath); err != nil {
		return common.NewHttpError("SBI rates file not found", http.StatusNotFound)
	}

	records, err := s.readCSVRecords(filePath)
	if err != nil {
		return err
	}

	rateMap, err := s.parseCSVToRateMap(records)
	if err != nil {
		return err
	}

	s.rateCache = rateMap
	return nil
}

func (s *SBIManagerImpl) readCSVRecords(filePath string) ([][]string, common.HttpError) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, common.NewServerError(err)
	}
	defer file.Close()

	// FIXME: #A Use Gocsv for parsing use SbiRate in Models
	reader := csv.NewReader(file)

	// Validate header
	header, err := reader.Read()
	if err != nil {
		return nil, common.NewServerError(err)
	}
	if len(header) < 3 || header[0] != "DATE" || header[1] != "TT BUY" || header[2] != "TT SELL" {
		return nil, common.NewHttpError("invalid CSV header format", http.StatusInternalServerError)
	}

	// Read remaining records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, common.NewServerError(err)
	}

	return records, nil
}
