package clients

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fa"
	"github.com/go-resty/resty/v2"
)

type SBIClient interface {
	FetchExchangeRates(ctx context.Context) (fa.ExchangeRates, common.HttpError)
}

type SBIClientImpl struct {
	baseUrl string
	client  *resty.Client
}

// BUG: Accept Base URL in Constructor.
func NewSBIClient(client *resty.Client) SBIClient {
	return &SBIClientImpl{
		baseUrl: "https://raw.githubusercontent.com/sahilgupta/sbi-fx-ratekeeper/main/csv_files/SBI_REFERENCE_RATES_USD.csv",
		client:  client,
	}
}

func (s *SBIClientImpl) FetchExchangeRates(ctx context.Context) (rates fa.ExchangeRates, err common.HttpError) {
	response, resErr := s.client.R().
		SetContext(ctx).
		Get(s.baseUrl)

	if err = util.ResponseProcessor(response, resErr); err != nil {
		return
	}

	// Parse CSV
	// FIXME: Break into Private Method and Custom Errors for Parsing
	reader := csv.NewReader(strings.NewReader(string(response.Body())))
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
