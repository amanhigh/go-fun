package manager

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
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

func (s *SBIManagerImpl) GetTTBuyRate(ctx context.Context, requestedDate time.Time) (rate float64, err common.HttpError) {
	rates, err := s.exchangeRepo.GetAllRecords(ctx)
	if err != nil {
		return 0, err
	}

	if len(rates) == 0 {
		return 0, common.ErrNotFound
	}

	// Try exact date match first
	dateStr := requestedDate.Format(time.DateOnly)
	for _, rate := range rates {
		if strings.Split(rate.Date, " ")[0] == dateStr {
			return rate.TTBuy, nil
		}
	}

	// Find closest previous date
	var closestDate string
	var closestRate float64
	for _, rate := range rates {
		rateDate := strings.Split(rate.Date, " ")[0]
		if rateDate <= dateStr && (closestDate == "" || rateDate > closestDate) {
			closestDate = rateDate
			closestRate = rate.TTBuy
		}
	}

	if closestDate != "" {
		parsedClosest, _ := time.Parse(time.DateOnly, closestDate)
		return closestRate, tax.NewClosestDateError(requestedDate, parsedClosest)
	}

	return 0, common.ErrNotFound
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
