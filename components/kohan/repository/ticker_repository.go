package repository

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// tickerSortColumn maps API sort-by field names to DB column names.
func tickerSortColumn(sortBy string) string {
	switch sortBy {
	case "ticker", "":
		return "external_id"
	default:
		return sortBy
	}
}

// TickerRepository provides persistence operations for barkat tickers.
type TickerRepository interface {
	util.BaseDbRepositoryInterface
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
	tx = tx.Order(clause.OrderByColumn{
		Column: clause.Column{Name: tickerSortColumn(query.SortBy)},
		Desc:   query.SortOrder == common.SortOrderDesc,
	})

	var tickers []barkat.Ticker
	err := tx.Offset(query.Offset).Limit(query.Limit).Find(&tickers).Error
	return tickers, err
}
