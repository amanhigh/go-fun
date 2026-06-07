package repository

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"gorm.io/gorm"
)

// PriceAlertRepository provides persistence operations for price alerts.
type PriceAlertRepository interface {
	util.BaseDbRepositoryInterface
	// GetByPairId resolves a pair id to exactly one AlertTicker.
	GetByPairId(ctx context.Context, pairID string) (barkat.AlertTicker, common.HttpError)
	// GetByTicker resolves the first AlertTicker under a parent ticker.
	GetByTicker(ctx context.Context, ticker string) (barkat.AlertTicker, common.HttpError)
	// ReplaceAlerts deletes existing alerts for owners and inserts replacement rows.
	ReplaceAlerts(ctx context.Context, alerts []barkat.PriceAlert) common.HttpError
	// DeleteByAlertID deletes one canonical alert by external alert id.
	DeleteByAlertID(ctx context.Context, alertID string) common.HttpError
	// ListPriceAlerts returns filtered, sorted, paginated price alerts.
	ListPriceAlerts(ctx context.Context, query barkat.PriceAlertQuery) ([]barkat.PriceAlert, int64, common.HttpError)
}

type PriceAlertRepositoryImpl struct {
	util.BaseDbRepository
}

var _ PriceAlertRepository = (*PriceAlertRepositoryImpl)(nil)

// NewPriceAlertRepository creates a new PriceAlertRepository backed by GORM.
func NewPriceAlertRepository(db *gorm.DB) *PriceAlertRepositoryImpl {
	return &PriceAlertRepositoryImpl{BaseDbRepository: util.NewBaseDbRepository(db)}
}

func (r *PriceAlertRepositoryImpl) GetByPairId(ctx context.Context, pairID string) (barkat.AlertTicker, common.HttpError) {
	var alertTicker barkat.AlertTicker
	if err := r.SafeTx(ctx).Where("pair_id = ?", pairID).First(&alertTicker).Error; err != nil {
		return barkat.AlertTicker{}, util.GormErrorMapper(err)
	}
	return alertTicker, nil
}

func (r *PriceAlertRepositoryImpl) GetByTicker(ctx context.Context, ticker string) (barkat.AlertTicker, common.HttpError) {
	var parent barkat.Ticker
	if httpErr := r.GetByExternalId(ctx, ticker, &parent); httpErr != nil {
		return barkat.AlertTicker{}, httpErr
	}

	var alertTicker barkat.AlertTicker
	if err := r.SafeTx(ctx).Where("ticker_id = ?", parent.ID).Order("id ASC").First(&alertTicker).Error; err != nil {
		return barkat.AlertTicker{}, util.GormErrorMapper(err)
	}
	return alertTicker, nil
}

func (r *PriceAlertRepositoryImpl) ReplaceAlerts(ctx context.Context, alerts []barkat.PriceAlert) common.HttpError {
	if len(alerts) == 0 {
		return nil
	}

	tx := r.SafeTx(ctx)

	alertTickerIDs := make([]uint64, 0, len(alerts))
	seen := make(map[uint64]bool, len(alerts))
	for _, a := range alerts {
		if !seen[a.AlertTickerID] {
			alertTickerIDs = append(alertTickerIDs, a.AlertTickerID)
			seen[a.AlertTickerID] = true
		}
	}

	if err := tx.Where("alert_ticker_id IN ?", alertTickerIDs).Delete(&barkat.PriceAlert{}).Error; err != nil {
		return util.GormErrorMapper(err)
	}
	if err := tx.Create(&alerts).Error; err != nil {
		return util.GormErrorMapper(err)
	}
	return nil
}

func (r *PriceAlertRepositoryImpl) DeleteByAlertID(ctx context.Context, alertID string) common.HttpError {
	result := r.SafeTx(ctx).Where("alert_id = ?", alertID).Delete(&barkat.PriceAlert{})
	if result.Error != nil {
		return util.GormErrorMapper(result.Error)
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound
	}
	return nil
}

func (r *PriceAlertRepositoryImpl) ListPriceAlerts(ctx context.Context, query barkat.PriceAlertQuery) ([]barkat.PriceAlert, int64, common.HttpError) {
	filteredTx := r.SafeTx(ctx).Model(&barkat.PriceAlert{})
	if query.Ticker != "" {
		filteredTx = filteredTx.Where("alert_ticker_id IN (SELECT at.id FROM alert_tickers at JOIN tickers t ON t.id = at.ticker_id WHERE t.external_id = ?)", query.Ticker)
	}

	var total int64
	if err := filteredTx.Count(&total).Error; err != nil {
		return nil, 0, util.GormErrorMapper(err)
	}

	sortedTx := util.ApplySort(filteredTx, util.SortOptions{
		SortBy:           query.SortBy,
		SortOrder:        query.SortOrder,
		DefaultSortBy:    "trigger_price",
		DefaultSortOrder: common.SortOrderAsc,
	})

	var alerts []barkat.PriceAlert
	if err := sortedTx.Preload("AlertTicker").
		Offset(query.Offset).Limit(query.Limit).Find(&alerts).Error; err != nil {
		return nil, 0, util.GormErrorMapper(err)
	}
	return alerts, total, nil
}
