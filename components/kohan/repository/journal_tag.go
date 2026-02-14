//nolint:dupl // Intentional CRUD pattern: Tag and Note repos follow same pattern for different sub-resources
package repository

import (
	"context"

	"github.com/amanhigh/go-fun/models/barkat"
	"gorm.io/gorm"
)

// TagRepository provides persistence operations for journal tags.
type TagRepository interface {
	// CreateTag attaches a new tag to an entry.
	CreateTag(ctx context.Context, tag *barkat.Tag) error
	// ListTags returns all tags for an entry, optionally filtered by type.
	ListTags(ctx context.Context, entryID string, tagType string) ([]barkat.Tag, error)
	// DeleteTag removes a tag by ID scoped to an entry.
	DeleteTag(ctx context.Context, entryID string, tagID string) error
}

type TagRepositoryImpl struct {
	db *gorm.DB
}

var _ TagRepository = (*TagRepositoryImpl)(nil)

// NewTagRepository creates a new TagRepository backed by GORM.
func NewTagRepository(db *gorm.DB) *TagRepositoryImpl {
	return &TagRepositoryImpl{db: db}
}

func (r *TagRepositoryImpl) CreateTag(ctx context.Context, tag *barkat.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *TagRepositoryImpl) ListTags(ctx context.Context, entryID, tagType string) ([]barkat.Tag, error) {
	var tags []barkat.Tag
	tx := r.db.WithContext(ctx).Where("entry_id = ?", entryID)
	if tagType != "" {
		tx = tx.Where("type = ?", tagType)
	}
	err := tx.Order("created_at").Find(&tags).Error
	return tags, err
}

func (r *TagRepositoryImpl) DeleteTag(ctx context.Context, entryID, tagID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND entry_id = ?", tagID, entryID).Delete(&barkat.Tag{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
