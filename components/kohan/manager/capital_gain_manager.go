package manager

import (
	"context"
	"net/http"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type CapitalGainManager interface {
	ProcessTaxGains(ctx context.Context, gains []tax.Gains) ([]tax.INRGains, common.HttpError)
}

type CapitalGainManagerImpl struct {
	sbiManager SBIManager
}

func NewCapitalGainManager(sbiManager SBIManager) *CapitalGainManagerImpl {
	return &CapitalGainManagerImpl{
		sbiManager: sbiManager,
	}
}

func (c *CapitalGainManagerImpl) ProcessTaxGains(_ context.Context, gains []tax.Gains) (taxGains []tax.INRGains, err common.HttpError) {
	for _, gain := range gains {
		var taxGain tax.INRGains
		// Copy base gains
		taxGain.Gains = gain

		// Parse sell date for exchange rate lookup
		var parseErr error
		if taxGain.TTDate, parseErr = gain.ParseSellDate(); parseErr != nil {
			return nil, common.NewHttpError(parseErr.Error(), http.StatusBadRequest)
		}

		// Get exchange rate for sell date
		if taxGain.TTRate, err = c.sbiManager.GetTTBuyRate(taxGain.TTDate); err != nil {
			return nil, err
		}

		taxGains = append(taxGains, taxGain)
	}
	return
}
