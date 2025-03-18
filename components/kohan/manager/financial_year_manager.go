package manager

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

//go:generate mockery --name FinancialYearManager
type FinancialYearManager[T tax.CSVRecord] interface {
	FilterRecordsByFY(ctx context.Context, records []T, year int) ([]T, common.HttpError)
}

type FinancialYearManagerImpl[T tax.CSVRecord] struct{}

func NewFinancialYearManager[T tax.CSVRecord]() FinancialYearManager[T] {
	return &FinancialYearManagerImpl[T]{}
}

func (f *FinancialYearManagerImpl[T]) FilterRecordsByFY(ctx context.Context, records []T, year int) ([]T, common.HttpError) {
	var filtered []T

	fyStart := time.Date(year, 4, 1, 0, 0, 0, 0, time.UTC)
	fyEnd := time.Date(year+1, 3, 31, 23, 59, 59, 0, time.UTC)

	for _, record := range records {
		date := record.GetDate()
		if (date.Equal(fyStart) || date.After(fyStart)) &&
			(date.Equal(fyEnd) || date.Before(fyEnd)) {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}
