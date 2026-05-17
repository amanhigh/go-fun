package repository

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"gorm.io/gorm"
)

// AlertTickerRepository provides persistence operations for Alert tickers.
type AlertTickerRepository interface {
	util.BaseDbRepositoryInterface
	// GetAlertTicker retrieves a single Alert ticker by symbol with parent ticker populated.
	GetAlertTicker(ctx context.Context, symbol string) (barkat.AlertTicker, common.HttpError)
	// ListAlertTickers returns a filtered, paginated list of Alert tickers with parent ticker populated.
	ListAlertTickers(ctx context.Context, query barkat.AlertTickerQuery) ([]barkat.AlertTicker, int64, common.HttpError)
}

type AlertTickerRepositoryImpl struct {
	util.BaseDbRepository
}

var _ AlertTickerRepository = (*AlertTickerRepositoryImpl)(nil)

// NewAlertTickerRepository creates a new AlertTickerRepository backed by GORM.
func NewAlertTickerRepository(db *gorm.DB) *AlertTickerRepositoryImpl {
	return &AlertTickerRepositoryImpl{
		BaseDbRepository: util.NewBaseDbRepository(db),
	}
}

// ---- Alert Ticker ----

func (r *AlertTickerRepositoryImpl) GetAlertTicker(ctx context.Context, symbol string) (barkat.AlertTicker, common.HttpError) {
	var result barkat.AlertTicker
	// FIXME: Remove Preload if Unused in UI.
	err := r.SafeTx(ctx).Model(&barkat.AlertTicker{}).
		Preload("ParentTicker").
		First(&result, "external_id = ?", symbol).Error
	return result, util.GormErrorMapper(err)
}

func (r *AlertTickerRepositoryImpl) ListAlertTickers(ctx context.Context, query barkat.AlertTickerQuery) ([]barkat.AlertTicker, int64, common.HttpError) {
	tx := r.applyAlertTickerFilters(r.SafeTx(ctx).Model(&barkat.AlertTicker{}), query)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, util.GormErrorMapper(err)
	}

	var alertTickers []barkat.AlertTicker
	if err := tx.Preload("ParentTicker").Offset(query.Offset).Limit(query.Limit).Find(&alertTickers).Error; err != nil {
		return nil, 0, util.GormErrorMapper(err)
	}

	return alertTickers, total, nil
}

func (r *AlertTickerRepositoryImpl) applyAlertTickerFilters(tx *gorm.DB, query barkat.AlertTickerQuery) *gorm.DB {
	if query.Symbol != "" {
		tx = tx.Where("external_id = ?", query.Symbol)
	}
	if query.Ticker != "" {
		tx = tx.Where("ticker_id IN (SELECT id FROM tickers WHERE external_id = ?)", query.Ticker)
	}
	if query.PairID != "" {
		tx = tx.Where("pair_id = ?", query.PairID)
	}
	if query.Exchange != "" {
		tx = tx.Where("exchange = ?", query.Exchange)
	}
	return tx
}
