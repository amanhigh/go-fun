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
	// GetByPairId resolves a pair id to an AlertTicker. When types are provided,
	// only the first type is matched (e.g. "PRIMARY"). When no type is given,
	// any AlertTicker for the pair id is returned.
	GetByPairId(ctx context.Context, pairID string, types ...string) (barkat.AlertTicker, common.HttpError)
	// GetByTicker resolves the first AlertTicker under a parent ticker.
	GetByTicker(ctx context.Context, ticker string) (barkat.AlertTicker, common.HttpError)
	// CreateAlerts inserts new price alert rows.
	CreateAlerts(ctx context.Context, alerts []barkat.PriceAlert) common.HttpError
	// DeleteByPairIDs deletes all price alerts for the given pair IDs.
	DeleteByPairIDs(ctx context.Context, pairIDs []string) common.HttpError
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

func (r *PriceAlertRepositoryImpl) GetByPairId(ctx context.Context, pairID string, types ...string) (barkat.AlertTicker, common.HttpError) {
	var alertTicker barkat.AlertTicker
	tx := r.SafeTx(ctx).Where(&barkat.AlertTicker{PairID: pairID})
	if len(types) > 0 {
		tx = tx.Where("type = ?", types[0])
	}
	if err := tx.First(&alertTicker).Error; err != nil {
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
	if err := r.SafeTx(ctx).Where(&barkat.AlertTicker{TickerID: parent.ID, Type: "PRIMARY"}).First(&alertTicker).Error; err != nil {
		return barkat.AlertTicker{}, util.GormErrorMapper(err)
	}
	return alertTicker, nil
}

func (r *PriceAlertRepositoryImpl) CreateAlerts(ctx context.Context, alerts []barkat.PriceAlert) common.HttpError {
	if len(alerts) == 0 {
		return nil
	}
	if err := r.SafeTx(ctx).Create(&alerts).Error; err != nil {
		return util.GormErrorMapper(err)
	}
	return nil
}

func (r *PriceAlertRepositoryImpl) DeleteByPairIDs(ctx context.Context, pairIDs []string) common.HttpError {
	if len(pairIDs) == 0 {
		return nil
	}
	tx := r.SafeTx(ctx)
	subQuery := tx.Model(&barkat.AlertTicker{}).Where("pair_id IN ?", pairIDs).Select("id")
	if err := tx.Where("alert_ticker_id IN (?)", subQuery).Delete(&barkat.PriceAlert{}).Error; err != nil {
		return util.GormErrorMapper(err)
	}
	return nil
}

func (r *PriceAlertRepositoryImpl) DeleteByAlertID(ctx context.Context, alertID string) common.HttpError {
	return r.DeleteBy(ctx, &barkat.PriceAlert{}, "alert_id = ?", alertID)
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
