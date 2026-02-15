package manager

import (
	"context"

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
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		return m.repo.Create(c, image)
	})
}

func (m *ImageManagerImpl) ListImages(ctx context.Context, entryID string) ([]barkat.Image, common.HttpError) {
	if httpErr := m.entryMgr.EntryExists(ctx, entryID); httpErr != nil {
		return nil, httpErr
	}

	var images []barkat.Image
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		var httpErr common.HttpError
		images, httpErr = m.repo.ListImages(c, entryID)
		return httpErr
	})
	return images, err
}

func (m *ImageManagerImpl) DeleteImage(ctx context.Context, entryID, imageID string) common.HttpError {
	if httpErr := m.entryMgr.EntryExists(ctx, entryID); httpErr != nil {
		return httpErr
	}
	image := &barkat.Image{EntryID: entryID}
	return m.repo.DeleteById(ctx, imageID, image)
}
