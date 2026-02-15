package repository

// ImageRepository provides persistence operations for journal images.
// Images are screenshots captured across multiple timeframes for each entry.

import (
	"context"
	"errors"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"gorm.io/gorm"
)

// ImageRepository provides persistence operations for journal images.
type ImageRepository interface {
	util.BaseDbRepositoryInterface
	// ListImages returns all images for an entry.
	ListImages(ctx context.Context, entryID string) ([]barkat.Image, common.HttpError)
}

type ImageRepositoryImpl struct {
	util.BaseDbRepository
}

var _ ImageRepository = (*ImageRepositoryImpl)(nil)

// NewImageRepository creates a new ImageRepository backed by GORM.
func NewImageRepository(db *gorm.DB) *ImageRepositoryImpl {
	return &ImageRepositoryImpl{BaseDbRepository: util.NewBaseDbRepository(db)}
}

func (r *ImageRepositoryImpl) ListImages(ctx context.Context, entryID string) ([]barkat.Image, common.HttpError) {
	var images []barkat.Image
	var txErr error
	if txErr = r.SafeTx(ctx).Where("entry_id = ?", entryID).Order("timeframe").Find(&images).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		return nil, util.GormErrorMapper(txErr)
	}
	return images, nil
}
