package clients

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/go-resty/resty/v2"
)

// go:generate mockery --name SBIClient
type SBIClient interface {
	FetchExchangeRates(ctx context.Context) (string, common.HttpError)
}

type SBIClientImpl struct {
	baseUrl string
	client  *resty.Client
}

func NewSBIClient(client *resty.Client, baseUrl string) *SBIClientImpl {
	return &SBIClientImpl{
		baseUrl: baseUrl,
		client:  client,
	}
}

func (s *SBIClientImpl) FetchExchangeRates(ctx context.Context) (result string, err common.HttpError) {
	response, resErr := s.client.R().
		SetResult(&result).
		SetContext(ctx).
		Get(s.baseUrl)

	err = util.ResponseProcessor(response, resErr)
	return
}
