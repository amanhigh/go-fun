package manager

import (
	"context"
	"sort"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type FinancialYearManager[T tax.CSVRecord] interface {
	FilterIndia(ctx context.Context, records []T, year int) ([]T, common.HttpError)
	FilterUS(ctx context.Context, records []T, year int) ([]T, common.HttpError)
}

type FinancialYearManagerImpl[T tax.CSVRecord] struct{}

func NewFinancialYearManager[T tax.CSVRecord]() FinancialYearManager[T] {
	return &FinancialYearManagerImpl[T]{}
}

func (f *FinancialYearManagerImpl[T]) FilterIndia(_ context.Context, records []T, year int) ([]T, common.HttpError) {
	var filtered []T

	fyStart := time.Date(year, 4, 1, 0, 0, 0, 0, time.UTC)
	fyEnd := time.Date(year+1, 3, 31, 23, 59, 59, 0, time.UTC)

	for _, record := range records {
		date, err := record.GetDate()
		if err != nil {
			return nil, err
		}
		if (date.Equal(fyStart) || date.After(fyStart)) &&
			(date.Equal(fyEnd) || date.Before(fyEnd)) {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}

func (f *FinancialYearManagerImpl[T]) FilterUS(_ context.Context, records []T, year int) ([]T, common.HttpError) {
	var filtered []T

	fyStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	fyEnd := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	// Sort Records based on date
	sort.SliceStable(records, func(i, j int) bool {
		// Ignore errors for sorting, handle them during filtering
		dateI, _ := records[i].GetDate()
		dateJ, _ := records[j].GetDate()
		return dateI.Before(dateJ)
	})

	for _, record := range records {
		date, err := record.GetDate()
		if err != nil {
			return nil, err
		}
		if (date.Equal(fyStart) || date.After(fyStart)) &&
			(date.Equal(fyEnd) || date.Before(fyEnd)) {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}
