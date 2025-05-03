package manager

import (
	"context"
	"os"
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
		return 0, tax.NewRateNotFoundError(requestedDate)
	}

	// Try exact match first
	if rate, found := s.findExactRate(rates, requestedDate); found {
		return rate, nil
	}

	// Find closest previous rate
	return s.findClosestRate(rates, requestedDate)
}

// findExactRate attempts to find exact date match
func (s *SBIManagerImpl) findExactRate(rates []tax.SbiRate, requestedDate time.Time) (rate float64, found bool) {
	// FIXME: #C Use exchangeRepo.GetRecordsForTicker which is now Date.
	dateStr := requestedDate.Format(time.DateOnly)
	for _, rate := range rates {
		rateDate, err := rate.GetDate()
		if err != nil {
			log.Error().Err(err).Str("Date", rate.Date).Msg("Failed to parse rate date")
			return 0, false
		}
		if rateDate.Format(time.DateOnly) == dateStr {
			return rate.TTBuy, true
		}
	}
	return 0, false
}

// findClosestRate finds closest previous rate and returns with ClosestDateError
func (s *SBIManagerImpl) findClosestRate(rates []tax.SbiRate, requestedDate time.Time) (float64, common.HttpError) {
	dateStr := requestedDate.Format(time.DateOnly)
	var closestDate time.Time
	var closestRate float64

	for _, rate := range rates {
		rateDate, err := rate.GetDate()
		if err != nil {
			return 0, err
		}
		rateDateStr := rateDate.Format(time.DateOnly)
		if rateDateStr <= dateStr && (closestDate.IsZero() || rateDate.After(closestDate)) {
			closestDate = rateDate
			closestRate = rate.TTBuy
		}
	}

	if !closestDate.IsZero() {
		return closestRate, tax.NewClosestDateError(requestedDate, closestDate)
	}

	return 0, tax.NewRateNotFoundError(requestedDate)
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
