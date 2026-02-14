package repository

import (
	"context"

	"github.com/amanhigh/go-fun/models/barkat"
	"gorm.io/gorm"
)

// JournalRepository provides persistence operations for journal entries and sub-resources.
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

	// CreateImage attaches a new image to an entry.
	CreateImage(ctx context.Context, image *barkat.Image) error
	// ListImages returns all images for an entry.
	ListImages(ctx context.Context, entryID string) ([]barkat.Image, error)
	// DeleteImage removes an image by ID scoped to an entry.
	DeleteImage(ctx context.Context, entryID string, imageID string) error

	// CreateNote attaches a new note to an entry.
	CreateNote(ctx context.Context, note *barkat.Note) error
	// ListNotes returns all notes for an entry, optionally filtered by status.
	ListNotes(ctx context.Context, entryID string, status string) ([]barkat.Note, error)
	// DeleteNote removes a note by ID scoped to an entry.
	DeleteNote(ctx context.Context, entryID string, noteID string) error

	// CreateTag attaches a new tag to an entry.
	CreateTag(ctx context.Context, tag *barkat.Tag) error
	// ListTags returns all tags for an entry, optionally filtered by type.
	ListTags(ctx context.Context, entryID string, tagType string) ([]barkat.Tag, error)
	// DeleteTag removes a tag by ID scoped to an entry.
	DeleteTag(ctx context.Context, entryID string, tagID string) error
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
	tx := r.db.WithContext(ctx).Model(&barkat.Entry{})

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

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := "created_at DESC"
	if query.SortBy != "" {
		direction := "DESC"
		if query.SortOrder == "asc" {
			direction = "ASC"
		}
		orderClause = query.SortBy + " " + direction
	}

	var entries []barkat.Entry
	err := tx.Order(orderClause).
		Offset(query.Offset).
		Limit(query.Limit).
		Find(&entries).Error

	return entries, total, err
}

func (r *JournalRepositoryImpl) EntryExists(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&barkat.Entry{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
