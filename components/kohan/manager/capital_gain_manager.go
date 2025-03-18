package manager

import (
	"context"

	repository "github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type CapitalGainManager interface {
	ProcessTaxGains(ctx context.Context, gains []tax.Gains) ([]tax.INRGains, common.HttpError)
	GetGainsForYear(ctx context.Context, year int) ([]tax.Gains, common.HttpError)
}

type CapitalGainManagerImpl struct {
	exchangeManager      ExchangeManager
	gainsRepository      repository.GainsRepository
	financialYearManager FinancialYearManager[tax.Gains]
}

func NewCapitalGainManager(exchangeManager ExchangeManager,
	gainsRepository repository.GainsRepository,
	financialYearManager FinancialYearManager[tax.Gains]) *CapitalGainManagerImpl {
	return &CapitalGainManagerImpl{
		exchangeManager:      exchangeManager,
		gainsRepository:      gainsRepository,
		financialYearManager: financialYearManager,
	}
}

func (c *CapitalGainManagerImpl) GetGainsForYear(ctx context.Context, year int) ([]tax.Gains, common.HttpError) {
	// Get all records from repository
	records, err := c.gainsRepository.GetAllRecords(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by financial year
	return c.financialYearManager.FilterRecordsByFY(ctx, records, year)
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
