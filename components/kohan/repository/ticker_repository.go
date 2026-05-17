package repository

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"gorm.io/gorm"
)

// TickerRepository provides persistence operations for barkat tickers.
type TickerRepository interface {
	util.BaseDbRepositoryInterface
	// GetTicker retrieves a single ticker by external_id with AlertTickers preloaded.
	GetTicker(ctx context.Context, ticker string) (barkat.Ticker, common.HttpError)
	// ListTickers returns a filtered, paginated list of tickers.
	ListTickers(ctx context.Context, query barkat.TickerQuery) ([]barkat.Ticker, int64, common.HttpError)
}

type TickerRepositoryImpl struct {
	util.BaseDbRepository
}

var _ TickerRepository = (*TickerRepositoryImpl)(nil)

// NewTickerRepository creates a new TickerRepository backed by GORM.
func NewTickerRepository(db *gorm.DB) *TickerRepositoryImpl {
	return &TickerRepositoryImpl{
		BaseDbRepository: util.NewBaseDbRepository(db),
	}
}

// ---- Ticker ----

func (r *TickerRepositoryImpl) GetTicker(ctx context.Context, ticker string) (barkat.Ticker, common.HttpError) {
	var result barkat.Ticker
	err := r.SafeTx(ctx).Model(&barkat.Ticker{}).
		Preload("AlertTickers").
		First(&result, &barkat.Ticker{Ticker: ticker}).Error
	return result, util.GormErrorMapper(err)
}

func (r *TickerRepositoryImpl) ListTickers(ctx context.Context, query barkat.TickerQuery) ([]barkat.Ticker, int64, common.HttpError) {
	tx := r.applyTickerFilters(r.SafeTx(ctx).Model(&barkat.Ticker{}), query)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, util.GormErrorMapper(err)
	}

	tickers, err := r.fetchTickers(tx, query)
	return tickers, total, util.GormErrorMapper(err)
}

func (r *TickerRepositoryImpl) applyTickerFilters(tx *gorm.DB, query barkat.TickerQuery) *gorm.DB {
	where := barkat.Ticker{}
	if query.Exchange != "" {
		exchange := query.Exchange
		where.Exchange = &exchange
	}
	if query.Type != "" {
		where.Type = query.Type
	}
	if query.State != "" {
		where.State = query.State
	}
	if query.Trend != "" {
		where.Trend = query.Trend
	}
	tx = tx.Where(&where)

	if query.Search != "" {
		like := "%" + query.Search + "%"
		tx = tx.Where("(external_id LIKE ? OR exchange LIKE ?)", like, like)
	}
	if query.IsFNO != nil {
		tx = tx.Where("is_fno = ?", *query.IsFNO)
	}
	if query.OpenedAfter != "" {
		tx = tx.Where("last_opened_at >= ?", query.OpenedAfter)
	}
	return tx
}

func (r *TickerRepositoryImpl) fetchTickers(tx *gorm.DB, query barkat.TickerQuery) ([]barkat.Ticker, error) {
	// Add alert_ticker_count via subquery
	tx = tx.Select("tickers.*, (SELECT count(*) FROM alert_tickers WHERE alert_tickers.ticker_id = tickers.id) AS alert_ticker_count")

	tx = util.ApplySort(tx, util.SortOptions{
		SortBy:           query.SortBy,
		SortOrder:        query.SortOrder,
		DefaultSortBy:    "ticker",
		DefaultSortOrder: common.SortOrderAsc,
		SortFieldMap: map[string]string{
			"ticker": "external_id",
		},
	})

	var tickers []barkat.Ticker
	err := tx.Offset(query.Offset).Limit(query.Limit).Find(&tickers).Error
	return tickers, err
}
