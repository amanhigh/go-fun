package repository

import (
	"context"

	"github.com/amanhigh/go-fun/models/barkat"
	"gorm.io/gorm"
)

func (r *JournalRepositoryImpl) CreateImage(ctx context.Context, image *barkat.Image) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *JournalRepositoryImpl) ListImages(ctx context.Context, entryID string) ([]barkat.Image, error) {
	var images []barkat.Image
	// TODO: Use Struct Based queries or these are fine ?
	err := r.db.WithContext(ctx).Where("entry_id = ?", entryID).Order("timeframe").Find(&images).Error
	return images, err
}

func (r *JournalRepositoryImpl) DeleteImage(ctx context.Context, entryID string, imageID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND entry_id = ?", imageID, entryID).Delete(&barkat.Image{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
