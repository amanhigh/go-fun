package repository

import (
	"context"

	"github.com/amanhigh/go-fun/models/barkat"
	"gorm.io/gorm"
)

// JournalRepository provides persistence operations for journal entries.
type JournalRepository interface {
	// CreateEntry persists a new journal entry with its images.
	CreateEntry(ctx context.Context, entry *barkat.Entry) error
	// GetEntry retrieves a single entry by ID with preloaded images.
	GetEntry(ctx context.Context, id string) (barkat.Entry, error)
	// ListEntries returns a filtered, paginated list of entries.
	ListEntries(ctx context.Context, query barkat.EntryQuery) ([]barkat.Entry, int64, error)
}

type JournalRepositoryImpl struct {
	db *gorm.DB
}

var _ JournalRepository = (*JournalRepositoryImpl)(nil)

// NewJournalRepository creates a new JournalRepository backed by GORM.
func NewJournalRepository(db *gorm.DB) *JournalRepositoryImpl {
	return &JournalRepositoryImpl{db: db}
}

func (r *JournalRepositoryImpl) CreateEntry(ctx context.Context, entry *barkat.Entry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *JournalRepositoryImpl) GetEntry(ctx context.Context, id string) (barkat.Entry, error) {
	var entry barkat.Entry
	err := r.db.WithContext(ctx).Preload("Images").First(&entry, "id = ?", id).Error
	return entry, err
}

func (r *JournalRepositoryImpl) ListEntries(ctx context.Context, query barkat.EntryQuery) ([]barkat.Entry, int64, error) {
	tx := r.db.WithContext(ctx).Model(&barkat.Entry{})

	if query.Ticker != "" {
		tx = tx.Where("ticker = ?", query.Ticker)
	}
	if query.Type != "" {
		tx = tx.Where("type = ?", query.Type)
	}
	if query.Outcome != "" {
		tx = tx.Where("outcome = ?", query.Outcome)
	}
	if query.Sequence != "" {
		tx = tx.Where("sequence = ?", query.Sequence)
	}
	if query.From != "" {
		tx = tx.Where("created_at >= ?", query.From)
	}
	if query.To != "" {
		tx = tx.Where("created_at <= ?", query.To)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var entries []barkat.Entry
	err := tx.Order("created_at DESC").
		Offset(query.Offset).
		Limit(query.Limit).
		Preload("Images").
		Find(&entries).Error

	return entries, total, err
}
