package manager

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fa"
)

type SBIManager interface {
	FetchAndParseExchangeRates(ctx context.Context) (fa.ExchangeRates, common.HttpError)
	DownloadRates(ctx context.Context) common.HttpError
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

func (s *SBIManagerImpl) FetchAndParseExchangeRates(ctx context.Context) (rates fa.ExchangeRates, err common.HttpError) {
	var csvContent string
	if csvContent, err = s.client.FetchExchangeRates(ctx); err != nil {
		return
	}

	// Parse CSV content
	reader := csv.NewReader(strings.NewReader(csvContent))
	return readRates(reader)
}

// SaveRates fetches the latest exchange rates and saves them to a CSV file
// in the configured download directory. Returns an error if any step fails.
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

func readRates(reader *csv.Reader) (rates fa.ExchangeRates, err common.HttpError) {
	// Skip header
	if _, err1 := reader.Read(); err1 != nil {
		err = common.NewServerError(err1)
		return
	}

	// Read rates
	for {
		record, err1 := reader.Read()
		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			err = common.NewServerError(err1)
			return
		}

		ttBuy, _ := strconv.ParseFloat(record[1], 64)
		ttSell, _ := strconv.ParseFloat(record[2], 64)

		rates.Rates = append(rates.Rates, fa.Rate{
			Date:   record[0],
			TTBuy:  ttBuy,
			TTSell: ttSell,
		})
	}
	return
}
