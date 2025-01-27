package clients

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	tax "github.com/amanhigh/go-fun/models/tax"
	"github.com/go-resty/resty/v2"
)

//go:generate mockery --name AlphaClient
type AlphaClient interface {
	FetchDailyPrices(ctx context.Context, ticker string) (tax.VantageStockData, common.HttpError)
}

type AlphaClientImpl struct {
	baseUrl string
	apiKey  string
	client  *resty.Client
}

func NewAlphaClient(client *resty.Client, baseURL string, apiKey string) *AlphaClientImpl {
	return &AlphaClientImpl{
		baseUrl: baseURL,
		apiKey:  apiKey,
		client:  client,
	}
}

func (a *AlphaClientImpl) FetchDailyPrices(ctx context.Context, ticker string) (stockData tax.VantageStockData, err common.HttpError) {
	response, resErr := a.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"function":   "TIME_SERIES_DAILY",
			"symbol":     ticker,
			"outputsize": "full",
			"apikey":     a.apiKey,
		}).
		SetResult(&stockData).
		Get(a.baseUrl)

	err = util.ResponseProcessor(response, resErr)
	return
}
