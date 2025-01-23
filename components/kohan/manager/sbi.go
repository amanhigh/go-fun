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
	"github.com/amanhigh/go-fun/models/fa"
)

func (s *SBIManagerImpl) readCSVRecords(filePath string) ([][]string, common.HttpError) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, common.NewServerError(err)
	}
	defer file.Close()

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

type SBIManager interface {
	DownloadRates(ctx context.Context) common.HttpError
	GetTTBuyRate(date time.Time) (float64, common.HttpError)
}

type SBIManagerImpl struct {
	client      clients.SBIClient
	downloadDir string
}

func NewSBIManager(client clients.SBIClient, downloadDir string) *SBIManagerImpl {
	return &SBIManagerImpl{
		client:      client,
		downloadDir: downloadDir,
	}
}

// SaveRates fetches the latest exchange rates and saves them to a CSV file
// in the configured download directory. Returns an error if any step fails.
func (s *SBIManagerImpl) GetTTBuyRate(date time.Time) (rate float64, err common.HttpError) {
	// Read CSV file from download directory
	filePath := filepath.Join(s.downloadDir, fa.SBI_RATES_FILENAME)

	// Check if file exists
	if _, err1 := os.Stat(filePath); err1 != nil {
		return 0, common.NewHttpError("SBI rates file not found", http.StatusNotFound)
	}

	// Read and validate CSV records
	records, err := s.readCSVRecords(filePath)
	if err != nil {
		return 0, err
	}

	// Format date in required format (YYYY-MM-DD)
	dateStr := date.Format("2006-01-02")

	// Search for matching date
	for _, record := range records {
		// Check if date matches (taking first part before space as done in Python)
		if strings.Split(record[0], " ")[0] == dateStr {
			rate, err1 := strconv.ParseFloat(record[1], 64)
			if err1 != nil {
				return 0, common.NewServerError(err1)
			}
			return rate, nil
		}
	}

	return 0, common.ErrNotFound
}

func (s *SBIManagerImpl) DownloadRates(ctx context.Context) (err common.HttpError) {
	var csvContent string
	if csvContent, err = s.client.FetchExchangeRates(ctx); err != nil {
		return
	}

	// Ensure directory exists
	if err1 := os.MkdirAll(s.downloadDir, os.ModePerm); err1 != nil {
		return common.NewServerError(err1)
	}

	// Write to file
	filePath := filepath.Join(s.downloadDir, fa.SBI_RATES_FILENAME)
	if err1 := os.WriteFile(filePath, []byte(csvContent), util.DEFAULT_PERM); err1 != nil {
		return common.NewServerError(err1)
	}

	return
}
