package manager

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

type FinancialYearManager interface {
	FilterRecordsByFY(ctx context.Context, records []tax.CSVRecord, year int) ([]tax.CSVRecord, common.HttpError)
}

type FinancialYearManagerImpl struct{}

func NewFinancialYearManager() FinancialYearManager {
	return &FinancialYearManagerImpl{}
}

func (f *FinancialYearManagerImpl) FilterRecordsByFY(ctx context.Context, records []tax.CSVRecord, year int) ([]tax.CSVRecord, common.HttpError) {
	var filtered []tax.CSVRecord

	// Financial year start and end
	fyStart := time.Date(year, 4, 1, 0, 0, 0, 0, time.UTC)
	fyEnd := time.Date(year+1, 3, 31, 23, 59, 59, 0, time.UTC)

	for _, record := range records {
		if date, err := record.GetDate(); err == nil {
			if (date.Equal(fyStart) || date.After(fyStart)) &&
				(date.Equal(fyEnd) || date.Before(fyEnd)) {
				filtered = append(filtered, record)
			}
		}
	}
	return filtered, nil
}
