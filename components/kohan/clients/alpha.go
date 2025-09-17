package clients

import (
	"context"
	"fmt"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	tax "github.com/amanhigh/go-fun/models/tax"
	"github.com/go-resty/resty/v2"
)

type AlphaClient interface {
	FetchDailyPrices(ctx context.Context, ticker string) (tax.VantageStockData, common.HttpError)
	ValidateAPIKey() common.HttpError
}

type AlphaClientImpl struct {
	baseUrl string
	apiKey  string
	client  *resty.Client
}

func NewAlphaClient(client *resty.Client, baseURL, apiKey string) *AlphaClientImpl {
	return &AlphaClientImpl{
		baseUrl: baseURL,
		apiKey:  apiKey,
		client:  client,
	}
}

func (a *AlphaClientImpl) ValidateAPIKey() common.HttpError {
	if strings.TrimSpace(a.apiKey) == "" {
		return common.NewServerError(fmt.Errorf("alpha Vantage API key is required for ticker download"))
	}
	return nil
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
