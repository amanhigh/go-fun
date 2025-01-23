package manager

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
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

	// Read and parse CSV file
	file, err1 := os.Open(filePath)
	if err1 != nil {
		return 0, common.NewServerError(err1)
	}
	defer file.Close()

	// Format date in required format (YYYY-MM-DD)
	dateStr := date.Format("2006-01-02")

	// Read CSV and find matching rate
	reader := csv.NewReader(file)
	// Skip header
	_, err1 = reader.Read()
	if err1 != nil {
		return 0, common.NewServerError(err1)
	}

	// Search for matching date
	for {
		record, err1 := reader.Read()
		if errors.Is(err1, io.EOF) {
			break
		}
		if err1 != nil {
			return 0, common.NewServerError(err1)
		}

		// Check if date matches (taking first part before space as done in Python)
		if strings.Split(record[0], " ")[0] == dateStr {
			rate, err1 = strconv.ParseFloat(record[1], 64)
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
