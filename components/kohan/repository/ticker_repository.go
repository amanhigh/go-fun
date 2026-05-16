package repository

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"gorm.io/gorm"
)

// tickerSortFieldMap maps API sort-by field names to DB column names.
var tickerSortFieldMap = map[string]string{
	"ticker":         "external_id",
	"exchange":       "exchange",
	"type":           "type",
	"state":          "state",
	"trend":          "trend",
	"last_opened_at": "last_opened_at",
}

// TickerRepository provides persistence operations for barkat tickers.
type TickerRepository interface {
	util.BaseDbRepositoryInterface
	// GetTicker retrieves a single ticker by its external_id with AlertTickers preloaded.
	GetTicker(ctx context.Context, externalId string) (barkat.Ticker, common.HttpError)
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

func (r *TickerRepositoryImpl) GetTicker(ctx context.Context, externalId string) (barkat.Ticker, common.HttpError) {
	var ticker barkat.Ticker
	err := r.SafeTx(ctx).Preload("AlertTickers").First(&ticker, "external_id = ?", externalId).Error
	return ticker, util.GormErrorMapper(err)
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
	if query.Search != "" {
		like := "%" + query.Search + "%"
		tx = tx.Where("(external_id LIKE ? OR exchange LIKE ?)", like, like)
	}
	if query.Exchange != "" {
		tx = tx.Where("exchange = ?", query.Exchange)
	}
	if query.Type != "" {
		tx = tx.Where("type = ?", query.Type)
	}
	if query.State != "" {
		tx = tx.Where("state = ?", query.State)
	}
	if query.Trend != "" {
		tx = tx.Where("trend = ?", query.Trend)
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
	orderClause := "external_id ASC"
	if query.SortBy != "" {
		direction := "ASC"
		if query.SortOrder == "desc" {
			direction = "DESC"
		}
		// Map API field name to DB column name (e.g. "ticker" → "external_id")
		colName := query.SortBy
		if mapped, ok := tickerSortFieldMap[query.SortBy]; ok {
			colName = mapped
		}
		orderClause = fmt.Sprintf("%s %s", colName, direction)
	}

	var tickers []barkat.Ticker
	err := tx.Order(orderClause).Offset(query.Offset).Limit(query.Limit).Find(&tickers).Error
	return tickers, err
}
