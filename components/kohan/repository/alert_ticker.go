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
	// GetAlertTicker retrieves a single Alert ticker by symbol with parent Ticker preloaded.
	GetAlertTicker(ctx context.Context, symbol string) (barkat.AlertTicker, common.HttpError)
	// ListAlertTickers returns a filtered, paginated list of Alert tickers with parent Ticker preloaded.
	ListAlertTickers(ctx context.Context, query barkat.AlertTickerQuery) ([]barkat.AlertTicker, int64, common.HttpError)
	// ExistsPrimaryAlertTicker checks if a PRIMARY Alert ticker already exists for the given ticker_id.
	ExistsPrimaryAlertTicker(ctx context.Context, tickerID uint64) (bool, common.HttpError)
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
	// Preload Ticker is required for AfterFind to populate TickerSymbol in API responses.
	err := r.SafeTx(ctx).Model(&barkat.AlertTicker{}).
		Preload("Ticker").
		First(&result, &barkat.AlertTicker{Symbol: symbol}).Error
	return result, util.GormErrorMapper(err)
}

func (r *AlertTickerRepositoryImpl) ListAlertTickers(ctx context.Context, query barkat.AlertTickerQuery) ([]barkat.AlertTicker, int64, common.HttpError) {
	tx := r.applyAlertTickerFilters(r.SafeTx(ctx).Model(&barkat.AlertTicker{}), query)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, util.GormErrorMapper(err)
	}

	var alertTickers []barkat.AlertTicker
	if err := tx.Preload("Ticker").Offset(query.Offset).Limit(query.Limit).Find(&alertTickers).Error; err != nil {
		return nil, 0, util.GormErrorMapper(err)
	}

	return alertTickers, total, nil
}

func (r *AlertTickerRepositoryImpl) applyAlertTickerFilters(tx *gorm.DB, query barkat.AlertTickerQuery) *gorm.DB {
	where := barkat.AlertTicker{}
	if query.Symbol != "" {
		where.Symbol = query.Symbol
	}
	if query.PairID != "" {
		where.PairID = query.PairID
	}
	if query.Exchange != "" {
		where.Exchange = &query.Exchange
	}
	if query.Type != "" {
		where.Type = query.Type
	}
	tx = tx.Where(&where)

	if query.Ticker != "" {
		tx = tx.Where("ticker_id IN (SELECT id FROM tickers WHERE external_id = ?)", query.Ticker)
	}
	return tx
}

func (r *AlertTickerRepositoryImpl) ExistsPrimaryAlertTicker(ctx context.Context, tickerID uint64) (bool, common.HttpError) {
	var count int64
	err := r.SafeTx(ctx).Model(&barkat.AlertTicker{}).
		Where(&barkat.AlertTicker{TickerID: tickerID, Type: barkat.AlertTickerTypePrimary}).
		Count(&count).Error
	if err != nil {
		return false, util.GormErrorMapper(err)
	}
	return count > 0, nil
}
