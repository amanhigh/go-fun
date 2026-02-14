package repository

import (
	"context"

	"github.com/amanhigh/go-fun/models/barkat"
	"gorm.io/gorm"
)

// JournalRepository provides persistence operations for journal entries.
//
//go:generate mockery --name JournalRepository
type JournalRepository interface {
	// CreateEntry persists a new journal entry with its associations.
	CreateEntry(ctx context.Context, entry *barkat.Entry) error
	// GetEntry retrieves a single entry by ID with preloaded associations.
	GetEntry(ctx context.Context, id string) (barkat.Entry, error)
	// ListEntries returns a filtered, paginated list of entry summaries (no associations).
	ListEntries(ctx context.Context, query barkat.EntryQuery) ([]barkat.Entry, int64, error)
	// EntryExists checks if an entry with the given ID exists.
	EntryExists(ctx context.Context, id string) (bool, error)
}

type JournalRepositoryImpl struct {
	db *gorm.DB
}

var _ JournalRepository = (*JournalRepositoryImpl)(nil)

// NewJournalRepository creates a new JournalRepository backed by GORM.
func NewJournalRepository(db *gorm.DB) *JournalRepositoryImpl {
	return &JournalRepositoryImpl{db: db}
}

// ---- Entry ----

func (r *JournalRepositoryImpl) CreateEntry(ctx context.Context, entry *barkat.Entry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *JournalRepositoryImpl) GetEntry(ctx context.Context, id string) (barkat.Entry, error) {
	var entry barkat.Entry
	err := r.db.WithContext(ctx).Preload("Images").Preload("Tags").Preload("Notes").First(&entry, "id = ?", id).Error
	return entry, err
}

func (r *JournalRepositoryImpl) ListEntries(ctx context.Context, query barkat.EntryQuery) ([]barkat.Entry, int64, error) {
	tx := r.applyEntryFilters(r.db.WithContext(ctx).Model(&barkat.Entry{}), query)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	entries, err := r.fetchEntries(tx, query)
	return entries, total, err
}

func (r *JournalRepositoryImpl) applyEntryFilters(tx *gorm.DB, query barkat.EntryQuery) *gorm.DB {
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

func (r *JournalRepositoryImpl) fetchEntries(tx *gorm.DB, query barkat.EntryQuery) ([]barkat.Entry, error) {
	orderClause := "created_at DESC"
	if query.SortBy != "" {
		direction := "DESC"
		if query.SortOrder == "asc" {
			direction = "ASC"
		}
		orderClause = query.SortBy + " " + direction
	}

	var entries []barkat.Entry
	err := tx.Order(orderClause).Offset(query.Offset).Limit(query.Limit).Find(&entries).Error
	return entries, err
}

func (r *JournalRepositoryImpl) EntryExists(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&barkat.Entry{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
