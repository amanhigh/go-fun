package manager

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// ImageManager provides business logic for journal image operations.
type ImageManager interface {
	// CreateImage attaches a new image to an entry.
	CreateImage(ctx context.Context, entryID string, image *barkat.Image) common.HttpError
	// ListImages returns all images for an entry.
	ListImages(ctx context.Context, entryID string) ([]barkat.Image, common.HttpError)
	// DeleteImage removes an image by ID scoped to an entry.
	DeleteImage(ctx context.Context, entryID string, imageID string) common.HttpError
}

type ImageManagerImpl struct {
	entryMgr JournalManager
	repo     repository.ImageRepository
}

var _ ImageManager = (*ImageManagerImpl)(nil)

// NewImageManager creates a new ImageManager.
func NewImageManager(entryMgr JournalManager, repo repository.ImageRepository) *ImageManagerImpl {
	return &ImageManagerImpl{entryMgr: entryMgr, repo: repo}
}

func (m *ImageManagerImpl) CreateImage(ctx context.Context, entryID string, image *barkat.Image) common.HttpError {
	if httpErr := m.entryMgr.EntryExists(ctx, entryID); httpErr != nil {
		return httpErr
	}
	image.EntryID = entryID
	if err := m.repo.CreateImage(ctx, image); err != nil {
		return common.NewServerError(fmt.Errorf("failed to create image: %w", err))
	}
	return nil
}

func (m *ImageManagerImpl) ListImages(ctx context.Context, entryID string) ([]barkat.Image, common.HttpError) {
	if httpErr := m.entryMgr.EntryExists(ctx, entryID); httpErr != nil {
		return nil, httpErr
	}
	images, err := m.repo.ListImages(ctx, entryID)
	if err != nil {
		return nil, common.NewServerError(fmt.Errorf("failed to list images: %w", err))
	}
	return images, nil
}

func (m *ImageManagerImpl) DeleteImage(ctx context.Context, entryID string, imageID string) common.HttpError {
	if err := m.repo.DeleteImage(ctx, entryID, imageID); err != nil {
		return common.ErrNotFound
	}
	return nil
}
