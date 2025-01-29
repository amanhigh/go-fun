package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type InterestManager interface {
	ProcessInterests(ctx context.Context, interests []tax.Interest) ([]tax.INRInterest, common.HttpError)
}

type InterestManagerImpl struct {
	exchangeManager ExchangeManager
}

func NewInterestManager(exchangeManager ExchangeManager) *InterestManagerImpl {
	return &InterestManagerImpl{
		exchangeManager: exchangeManager,
	}
}

func (i *InterestManagerImpl) ProcessInterests(ctx context.Context, interests []tax.Interest) (inrInterests []tax.INRInterest, err common.HttpError) {
	exchangeables := make([]tax.Exchangeable, 0, len(interests))
	for _, interest := range interests {
		var inrInterest tax.INRInterest
		inrInterest.Interest = interest

		inrInterests = append(inrInterests, inrInterest)
		exchangeables = append(exchangeables, &inrInterest)
	}

	err = i.exchangeManager.Exchange(ctx, exchangeables)
	return
}
