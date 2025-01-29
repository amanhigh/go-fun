package manager

import (
	"context"
	"os"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/rs/zerolog/log"
)

//go:generate mockery --name SBIManager
type SBIManager interface {
	DownloadRates(ctx context.Context) common.HttpError
	GetTTBuyRate(ctx context.Context, date time.Time) (float64, common.HttpError)
}

type SBIManagerImpl struct {
	client       clients.SBIClient
	filePath     string
	exchangeRepo repository.ExchangeRepository
}

func NewSBIManager(client clients.SBIClient, filePath string, exchangeRepo repository.ExchangeRepository) *SBIManagerImpl {
	return &SBIManagerImpl{
		client:       client,
		filePath:     filePath,
		exchangeRepo: exchangeRepo,
	}
}

func (s *SBIManagerImpl) GetTTBuyRate(ctx context.Context, date time.Time) (rate float64, err common.HttpError) {
	// Get rates for date using repository
	rates, err := s.exchangeRepo.GetAllRecords(ctx, date)
	if err != nil {
		return 0, err
	}

	if len(rates) == 0 {
		return 0, common.ErrNotFound
	}

	return rates[0].TTBuy, nil
}

func (s *SBIManagerImpl) DownloadRates(ctx context.Context) (err common.HttpError) {
	// Skip if file exists
	if _, err := os.Stat(s.filePath); err == nil {
		log.Info().Str("Path", s.filePath).Msg("SBI rates file already exists, skipping download")
		return nil
	}

	var csvContent string
	if csvContent, err = s.client.FetchExchangeRates(ctx); err != nil {
		return
	}

	// Write to file
	if err1 := os.WriteFile(s.filePath, []byte(csvContent), util.DEFAULT_PERM); err1 != nil {
		return common.NewServerError(err1)
	}

	return
}
