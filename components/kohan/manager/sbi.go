package manager

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fa"
)

type SBIManager interface {
	FetchAndParseExchangeRates(ctx context.Context) (fa.ExchangeRates, common.HttpError)
}

type SBIManagerImpl struct {
	client clients.SBIClient
}

func NewSBIManager(client clients.SBIClient) *SBIManagerImpl {
	return &SBIManagerImpl{client: client}
}

func (s *SBIManagerImpl) FetchAndParseExchangeRates(ctx context.Context) (rates fa.ExchangeRates, err common.HttpError) {
	var csvContent string
	if csvContent, err = s.client.FetchExchangeRates(ctx); err != nil {
		return
	}

	// Parse CSV content
	reader := csv.NewReader(strings.NewReader(csvContent))
	// Skip header
	if _, err1 := reader.Read(); err1 != nil {
		err = common.NewHttpError(err1.Error(), 500)
		return
	}

	// Read rates
	for {
		record, err1 := reader.Read()
		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			err = common.NewHttpError(err1.Error(), 500)
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
