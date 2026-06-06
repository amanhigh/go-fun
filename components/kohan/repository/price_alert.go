package repository

import (
	"context"
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"gorm.io/gorm"
)

// PriceAlertRepository provides persistence operations for price alerts.
type PriceAlertRepository interface {
	util.BaseDbRepositoryInterface
	// ResolveAlertTickerByPairID resolves a pair id to exactly one AlertTicker.
	ResolveAlertTickerByPairID(ctx context.Context, pairID string) (barkat.AlertTicker, common.HttpError)
	// GetFirstAlertTickerForTicker resolves the first AlertTicker under a parent ticker.
	GetFirstAlertTickerForTicker(ctx context.Context, ticker string) (barkat.AlertTicker, common.HttpError)
	// ReplaceAlerts deletes existing alerts for owners and inserts replacement rows.
	ReplaceAlerts(ctx context.Context, alertTickerIDs []uint64, alerts []barkat.PriceAlert) common.HttpError
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

func (r *PriceAlertRepositoryImpl) ResolveAlertTickerByPairID(ctx context.Context, pairID string) (barkat.AlertTicker, common.HttpError) {
	var alertTickers []barkat.AlertTicker
	if err := r.SafeTx(ctx).Where("pair_id = ?", pairID).Find(&alertTickers).Error; err != nil {
		return barkat.AlertTicker{}, util.GormErrorMapper(err)
	}
	switch len(alertTickers) {
	case 0:
		return barkat.AlertTicker{}, common.ErrNotFound
	case 1:
		return alertTickers[0], nil
	default:
		return barkat.AlertTicker{}, common.NewHttpError("Ambiguous pair id ownership", http.StatusConflict)
	}
}

func (r *PriceAlertRepositoryImpl) GetFirstAlertTickerForTicker(ctx context.Context, ticker string) (barkat.AlertTicker, common.HttpError) {
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

func (r *PriceAlertRepositoryImpl) ReplaceAlerts(ctx context.Context, alertTickerIDs []uint64, alerts []barkat.PriceAlert) common.HttpError {
	if len(alertTickerIDs) == 0 {
		return nil
	}
	if err := r.SafeTx(ctx).Where("alert_ticker_id IN ?", alertTickerIDs).Delete(&barkat.PriceAlert{}).Error; err != nil {
		return util.GormErrorMapper(err)
	}
	if len(alerts) == 0 {
		return nil
	}
	if err := r.SafeTx(ctx).Create(&alerts).Error; err != nil {
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
	tx := r.applyPriceAlertFilters(r.SafeTx(ctx).Model(&barkat.PriceAlert{}), query)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, util.GormErrorMapper(err)
	}

	tx = util.ApplySort(tx, util.SortOptions{
		SortBy:           query.SortBy,
		SortOrder:        query.SortOrder,
		DefaultSortBy:    "trigger_price",
		DefaultSortOrder: common.SortOrderAsc,
	})

	var alerts []barkat.PriceAlert
	if err := tx.Preload("AlertTicker").
		Offset(query.Offset).Limit(query.Limit).Find(&alerts).Error; err != nil {
		return nil, 0, util.GormErrorMapper(err)
	}
	return alerts, total, nil
}

func (r *PriceAlertRepositoryImpl) applyPriceAlertFilters(tx *gorm.DB, query barkat.PriceAlertQuery) *gorm.DB {
	if query.Ticker != "" {
		tx = tx.Where("alert_ticker_id IN (SELECT at.id FROM alert_tickers at JOIN tickers t ON t.id = at.ticker_id WHERE t.external_id = ?)", query.Ticker)
	}
	return tx
}
