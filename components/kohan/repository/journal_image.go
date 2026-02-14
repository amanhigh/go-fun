package repository

// ImageRepository provides persistence operations for journal images.
// Images are screenshots captured across multiple timeframes for each entry.

import (
	"context"

	"github.com/amanhigh/go-fun/models/barkat"
	"gorm.io/gorm"
)

// ImageRepository provides persistence operations for journal images.
type ImageRepository interface {
	// CreateImage attaches a new image to an entry.
	CreateImage(ctx context.Context, image *barkat.Image) error
	// ListImages returns all images for an entry.
	ListImages(ctx context.Context, entryID string) ([]barkat.Image, error)
	// DeleteImage removes an image by ID scoped to an entry.
	DeleteImage(ctx context.Context, entryID string, imageID string) error
}

type ImageRepositoryImpl struct {
	db *gorm.DB
}

var _ ImageRepository = (*ImageRepositoryImpl)(nil)

// NewImageRepository creates a new ImageRepository backed by GORM.
func NewImageRepository(db *gorm.DB) *ImageRepositoryImpl {
	return &ImageRepositoryImpl{db: db}
}

func (r *ImageRepositoryImpl) CreateImage(ctx context.Context, image *barkat.Image) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *ImageRepositoryImpl) ListImages(ctx context.Context, entryID string) ([]barkat.Image, error) {
	var images []barkat.Image
	// TODO: Use Struct Based queries or these are fine ?
	err := r.db.WithContext(ctx).Where("entry_id = ?", entryID).Order("timeframe").Find(&images).Error
	return images, err
}

func (r *ImageRepositoryImpl) DeleteImage(ctx context.Context, entryID, imageID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND entry_id = ?", imageID, entryID).Delete(&barkat.Image{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
