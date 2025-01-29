package manager

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type CapitalGainManager interface {
	ProcessTaxGains(ctx context.Context, gains []tax.Gains) ([]tax.INRGains, common.HttpError)
	// FIXME: #A Expose GetAllGains and Ingrate wit Tax Manager ?
}

type CapitalGainManagerImpl struct {
	exchangeManager ExchangeManager
}

func NewCapitalGainManager(exchangeManager ExchangeManager) *CapitalGainManagerImpl {
	return &CapitalGainManagerImpl{
		exchangeManager: exchangeManager,
	}
}

func (c *CapitalGainManagerImpl) ProcessTaxGains(ctx context.Context, gains []tax.Gains) (taxGains []tax.INRGains, err common.HttpError) {
	exchangeableGains := make([]tax.Exchangeable, len(taxGains))
	for _, gain := range gains {
		var taxGain tax.INRGains
		// Copy base gains
		taxGain.Gains = gain

		taxGains = append(taxGains, taxGain)
		exchangeableGains = append(exchangeableGains, &taxGain)
	}
	err = c.exchangeManager.Exchange(ctx, exchangeableGains)
	return
}
