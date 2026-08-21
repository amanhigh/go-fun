package repository

// ImageRepository provides persistence operations for journal images.
// Images are ordered by date first, then higher-to-lower timeframe.
// Images are screenshots captured across multiple timeframes for each journal.

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
	util.BaseDbRepository
	// ListImages returns all images for a journal, optionally filtered by image type.
	ListImages(ctx context.Context, journalID uint64, imageType string) ([]barkat.Image, common.HttpError)
}

type ImageRepositoryImpl struct {
	util.BaseDbRepository
}

var _ ImageRepository = (*ImageRepositoryImpl)(nil)

// NewImageRepository creates a new ImageRepository backed by GORM.
func NewImageRepository(baseRepository util.BaseDbRepository) *ImageRepositoryImpl {
	return &ImageRepositoryImpl{BaseDbRepository: baseRepository}
}

func (r *ImageRepositoryImpl) ListImages(ctx context.Context, journalID uint64, imageType string) ([]barkat.Image, common.HttpError) {
	var images []barkat.Image
	var txErr error

	where := barkat.Image{JournalID: journalID}
	if imageType != "" {
		where.ImageType = imageType
	}
	query := r.SafeTx(ctx).Where(&where)

	if txErr = query.
		Order("DATE(created_at) ASC").
		Order(ImageTimeframeOrder + " DESC").
		Find(&images).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		return nil, util.GormErrorMapper(txErr)
	}
	return images, nil
}
