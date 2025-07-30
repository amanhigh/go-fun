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
	// TODO: Get Last TT Buy Rate for month.
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
	// Try exact match first using repository's direct lookup
	rate, err = s.findExactRate(ctx, requestedDate)
	if err == nil {
		return rate, nil
	}

	// If the error is RateNotFoundError, try finding the closest rate by fetching all records.
	// Otherwise, return the original error (e.g., ServerError from date parsing).
	if _, ok := err.(tax.RateNotFoundError); ok {
		rates, repoErr := s.exchangeRepo.GetAllRecords(ctx)
		if repoErr != nil {
			return 0, common.NewServerError(repoErr)
		}

		if len(rates) == 0 {
			// If no records at all, still a RateNotFoundError
			return 0, tax.NewRateNotFoundError(requestedDate)
		}

		// Find closest previous rate from all records
		return s.findClosestRate(rates, requestedDate)
	}

	// Return the error if it's not a RateNotFoundError
	return 0, err
}

// findExactRate attempts to find exact date match using the repository's direct lookup.
func (s *SBIManagerImpl) findExactRate(ctx context.Context, requestedDate time.Time) (rate float64, err common.HttpError) {
	dateStr := requestedDate.Format(time.DateOnly)
	rateRecords, repoErr := s.exchangeRepo.GetRecordsForTicker(ctx, dateStr)
	if repoErr != nil {
		return 0, common.NewServerError(repoErr)
	}
	if len(rateRecords) > 0 {
		return rateRecords[0].TTBuy, nil
	}
	return 0, tax.NewRateNotFoundError(requestedDate)
}

// findClosestRate finds closest previous rate and returns with ClosestDateError
func (s *SBIManagerImpl) findClosestRate(rates []tax.SbiRate, requestedDate time.Time) (float64, common.HttpError) {
	dateStr := requestedDate.Format(time.DateOnly)
	var closestDate time.Time
	var closestRate float64

	for _, rate := range rates {
		rateDate, dateErr := rate.GetDate()
		if dateErr != nil {
			return 0, dateErr
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
	var fetchErr error // Rename error variable to avoid shadowing
	if csvContent, fetchErr = s.client.FetchExchangeRates(ctx); fetchErr != nil {
		return common.NewServerError(fetchErr) // Wrap the standard error
	}

	// Write to file
	if err1 := os.WriteFile(s.filePath, []byte(csvContent), util.DEFAULT_PERM); err1 != nil {
		return common.NewServerError(err1)
	}

	return nil
}
