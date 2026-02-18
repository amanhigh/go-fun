package manager

// ImageManager provides business logic for journal image operations.
// Images represent visual attachments to journal entries with timeframe metadata.

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

type ImageManager interface {
	// CreateImage attaches a new image to an entry.
	CreateImage(ctx context.Context, entryID string, image barkat.Image) (*barkat.Image, common.HttpError)
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

func (m *ImageManagerImpl) CreateImage(ctx context.Context, entryID string, image barkat.Image) (*barkat.Image, common.HttpError) {
	image.EntryID = entryID
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Check entry existence within transaction
		if httpErr := m.entryMgr.EntryExists(c, entryID); httpErr != nil {
			return httpErr
		}
		return m.repo.Create(c, &image)
	})
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (m *ImageManagerImpl) ListImages(ctx context.Context, entryID string) ([]barkat.Image, common.HttpError) {
	var images []barkat.Image
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Check entry existence within transaction
		if httpErr := m.entryMgr.EntryExists(c, entryID); httpErr != nil {
			return httpErr
		}
		var httpErr common.HttpError
		images, httpErr = m.repo.ListImages(c, entryID)
		return httpErr
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (m *ImageManagerImpl) DeleteImage(ctx context.Context, entryID, imageID string) common.HttpError {
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Check entry existence within transaction
		if httpErr := m.entryMgr.EntryExists(c, entryID); httpErr != nil {
			return httpErr
		}
		return m.repo.DeleteById(c, imageID, &barkat.Image{EntryID: entryID})
	})
}
