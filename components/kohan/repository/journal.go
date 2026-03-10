package repository

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"gorm.io/gorm"
)

// JournalRepository provides persistence operations for journals.
type JournalRepository interface {
	util.BaseDbRepositoryInterface
	// GetJournal retrieves a single journal by ID with preloaded associations.
	GetJournal(ctx context.Context, id string) (barkat.Journal, error)
	// ListJournals returns a filtered, paginated list of journal summaries (no associations).
	ListJournals(ctx context.Context, query barkat.JournalQuery) ([]barkat.Journal, int64, error)
}

type JournalRepositoryImpl struct {
	util.BaseDbRepository
}

var _ JournalRepository = (*JournalRepositoryImpl)(nil)

// NewJournalRepository creates a new JournalRepository backed by GORM.
func NewJournalRepository(db *gorm.DB) *JournalRepositoryImpl {
	return &JournalRepositoryImpl{
		BaseDbRepository: util.NewBaseDbRepository(db),
	}
}

// ---- Journal ----

func (r *JournalRepositoryImpl) GetJournal(ctx context.Context, id string) (barkat.Journal, error) {
	var journal barkat.Journal
	err := r.SafeTx(ctx).Preload("Images").Preload("Tags").Preload("Notes").First(&journal, "id = ?", id).Error
	return journal, err
}

func (r *JournalRepositoryImpl) ListJournals(ctx context.Context, query barkat.JournalQuery) ([]barkat.Journal, int64, error) {
	tx := r.applyJournalFilters(r.SafeTx(ctx).Model(&barkat.Journal{}), query)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	journals, err := r.fetchJournals(tx, query)
	return journals, total, err
}

func (r *JournalRepositoryImpl) applyJournalFilters(tx *gorm.DB, query barkat.JournalQuery) *gorm.DB {
	if query.Ticker != "" {
		tx = tx.Where("ticker = ?", query.Ticker)
	}
	if query.Type != "" {
		tx = tx.Where("type = ?", query.Type)
	}
	if query.Status != "" {
		tx = tx.Where("status = ?", query.Status)
	}
	if query.Sequence != "" {
		tx = tx.Where("sequence = ?", query.Sequence)
	}
	if query.CreatedAfter != "" {
		tx = tx.Where("created_at >= ?", query.CreatedAfter)
	}
	if query.CreatedBefore != "" {
		tx = tx.Where("created_at <= ?", query.CreatedBefore)
	}
	return tx
}

func (r *JournalRepositoryImpl) fetchJournals(tx *gorm.DB, query barkat.JournalQuery) ([]barkat.Journal, error) {
	orderClause := "created_at DESC"
	if query.SortBy != "" {
		direction := "DESC"
		if query.SortOrder == "asc" {
			direction = "ASC"
		}
		orderClause = query.SortBy + " " + direction
	}

	var journals []barkat.Journal
	err := tx.Preload("Images").Order(orderClause).Offset(query.Offset).Limit(query.Limit).Find(&journals).Error
	return journals, err
}
