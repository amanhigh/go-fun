package manager

import (
	"context"
	"fmt"
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
	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
)

//go:generate mockery --name SBIManager
type SBIManager interface {
	DownloadRates(ctx context.Context) common.HttpError
	GetTTBuyRate(date time.Time) (float64, common.HttpError)
}

type SBIManagerImpl struct {
	client clients.SBIClient
	// BUG: Change to File Path done in provider
	downloadDir string
	rateCache   map[string]float64
}

func NewSBIManager(client clients.SBIClient, downloadDir string) *SBIManagerImpl {
	return &SBIManagerImpl{
		client:      client,
		downloadDir: downloadDir,
		// FIXME: Link to Exchange Repo Clean Logic and Test Retain Caching Logic.
		rateCache: make(map[string]float64),
	}
}

// ratesFileExists checks if the SBI rates file already exists
func (s *SBIManagerImpl) ratesFileExists() bool {
	filePath := filepath.Join(s.downloadDir, tax.SBI_RATES_FILENAME)
	_, err := os.Stat(filePath)
	return err == nil
}

// SaveRates fetches the latest exchange rates and saves them to a CSV file
// in the configured download directory. Returns an error if any step fails.
func (s *SBIManagerImpl) GetTTBuyRate(date time.Time) (rate float64, err common.HttpError) {
	if err = s.loadRatesIfNeeded(); err != nil {
		return 0, err
	}

	dateStr := date.Format(time.DateOnly)
	if rate, exists := s.rateCache[dateStr]; exists {
		return rate, nil
	}
	return 0, common.ErrNotFound
}

func (s *SBIManagerImpl) DownloadRates(ctx context.Context) (err common.HttpError) {
	// Skip if file exists
	if s.ratesFileExists() {
		log.Info().Str("Path", s.downloadDir).Msg("SBI rates file already exists, skipping download")
		return nil
	}

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

	var rates []tax.SbiRate
	// TODO: Move to BaseCSVRepository
	if err := gocsv.Unmarshal(file, &rates); err != nil {
		return nil, common.NewServerError(err)
	}

	// Validate using model's IsValid method
	if len(rates) == 0 || !rates[0].IsValid() {
		return nil, common.NewHttpError("invalid CSV header format", http.StatusInternalServerError)
	}

	// Convert to string slice format for existing code
	var records [][]string
	for _, rate := range rates {
		records = append(records, []string{rate.Date, fmt.Sprintf("%.2f", rate.TTBuy), fmt.Sprintf("%.2f", rate.TTSell)})
	}

	return records, nil
}
