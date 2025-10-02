package manager

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/rs/zerolog/log"
)

type SBIManager interface {
	DownloadRates(ctx context.Context) common.HttpError
	GetTTBuyRate(ctx context.Context, date time.Time) (float64, common.HttpError)
	// GetLastMonthEndRate returns the last available TT Buy rate for the month preceding the given date.
	// It precomputes and caches all month-end rates on first call for performance.
	GetLastMonthEndRate(ctx context.Context, date time.Time) (tax.MonthEndRate, common.HttpError)
}

type SBIManagerImpl struct {
	client        clients.SBIClient
	filePath      string
	exchangeRepo  repository.ExchangeRepository
	monthEndCache map[string]tax.MonthEndRate
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
		if rateDateStr <= dateStr && rate.TTBuy > 0 && (closestDate.IsZero() || rateDate.After(closestDate)) {
			closestDate = rateDate
			closestRate = rate.TTBuy
		}
	}

	if !closestDate.IsZero() {
		return closestRate, tax.NewClosestDateError(requestedDate, closestDate)
	}

	return 0, tax.NewRateNotFoundError(requestedDate)
}

// GetLastMonthEndRate returns the last available TT Buy rate for the month
// immediately preceding the month of the given date.
//
// On first call, it precomputes and caches all month-end rates for performance.
// The cache is immutable after initialization, so no mutex is needed.
//
// Example: date=2024-02-20 → returns last rate from Jan 2024
func (s *SBIManagerImpl) GetLastMonthEndRate(ctx context.Context, date time.Time) (tax.MonthEndRate, common.HttpError) {
	if s.monthEndCache == nil {
		if err := s.buildMonthEndCache(ctx); err != nil {
			return tax.MonthEndRate{}, err
		}
	}

	// Calculate preceding month
	precedingMonth := getPrecedingMonth(date)
	cacheKey := fmt.Sprintf("%d-%02d", precedingMonth.Year(), precedingMonth.Month())

	// Lookup in cache (no lock needed - immutable)
	if cached, found := s.monthEndCache[cacheKey]; found {
		log.Debug().
			Time("InputDate", date).
			Time("PrecedingMonth", precedingMonth).
			Float64("Rate", cached.Rate).
			Time("ActualDate", cached.ActualDate).
			Msg("SBIManager: Month-end rate from cache")
		return cached, nil
	}

	return tax.MonthEndRate{}, tax.NewRateNotFoundError(precedingMonth)
}

// getPrecedingMonth calculates the first day of the month immediately preceding the given date's month
func getPrecedingMonth(date time.Time) time.Time {
	firstOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	return firstOfMonth.AddDate(0, -1, 0)
}

// buildMonthEndCache precomputes all month-end rates by grouping rates by month
// and keeping the latest date in each month
func (s *SBIManagerImpl) buildMonthEndCache(ctx context.Context) common.HttpError {
	rates, repoErr := s.exchangeRepo.GetAllRecords(ctx)
	if repoErr != nil {
		return common.NewServerError(repoErr)
	}

	if len(rates) == 0 {
		log.Warn().Msg("SBIManager: No rates available to build cache")
		s.monthEndCache = make(map[string]tax.MonthEndRate)
		return nil
	}

	monthMap := make(map[string]tax.MonthEndRate)

	for _, rate := range rates {
		rateDate, dateErr := rate.GetDate()
		if dateErr != nil {
			return dateErr
		}

		monthKey := fmt.Sprintf("%d-%02d", rateDate.Year(), rateDate.Month())

		// Keep the latest date in each month
		if existing, found := monthMap[monthKey]; !found || rateDate.After(existing.ActualDate) {
			monthMap[monthKey] = tax.MonthEndRate{
				Rate:       rate.TTBuy,
				ActualDate: rateDate,
			}
		}
	}

	s.monthEndCache = monthMap

	log.Info().
		Int("MonthsLoaded", len(monthMap)).
		Msg("SBIManager: Precomputed and cached all month-end rates")

	return nil
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
