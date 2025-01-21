package manager

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
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
	SaveRates(ctx context.Context) common.HttpError
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
func (s *SBIManagerImpl) SaveRates(ctx context.Context) (err common.HttpError) {
	var rates fa.ExchangeRates
	if rates, err = s.FetchAndParseExchangeRates(ctx); err == nil {
		err = s.saveRatesToFile(rates)
	}
	return
}

func (s *SBIManagerImpl) saveRatesToFile(rates fa.ExchangeRates) (err common.HttpError) {
	// BUG: Make file name constant.
	filePath := filepath.Join(s.downloadDir, "SBI_REFERENCE_RATES_USD.csv")

	// Create CSV content
	var content []string
	content = append(content, "DATE,TT BUY,TT SELL")
	for _, rate := range rates.Rates {
		line := fmt.Sprintf("%s,%.2f,%.2f", rate.Date, rate.TTBuy, rate.TTSell)
		content = append(content, line)
	}

	// Write file
	if writeErr := util.WriteLines(filePath, content); writeErr != nil {
		err = common.NewServerError(writeErr)
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
