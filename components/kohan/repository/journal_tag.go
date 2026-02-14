package repository

import (
	"context"

	"github.com/amanhigh/go-fun/models/barkat"
	"gorm.io/gorm"
)

func (r *JournalRepositoryImpl) CreateTag(ctx context.Context, tag *barkat.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *JournalRepositoryImpl) ListTags(ctx context.Context, entryID string, tagType string) ([]barkat.Tag, error) {
	var tags []barkat.Tag
	tx := r.db.WithContext(ctx).Where("entry_id = ?", entryID)
	if tagType != "" {
		tx = tx.Where("type = ?", tagType)
	}
	err := tx.Order("created_at").Find(&tags).Error
	return tags, err
}

func (r *JournalRepositoryImpl) DeleteTag(ctx context.Context, entryID string, tagID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND entry_id = ?", tagID, entryID).Delete(&barkat.Tag{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
