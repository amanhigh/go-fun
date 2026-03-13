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
	// CreateImage attaches a new image to a journal.
	CreateImage(ctx context.Context, journalID string, image barkat.Image) (*barkat.Image, common.HttpError)
	// ListImages returns all images for a journal.
	ListImages(ctx context.Context, journalID string) ([]barkat.Image, common.HttpError)
	// DeleteImage removes an image by ID scoped to a journal.
	DeleteImage(ctx context.Context, journalID string, imageID string) common.HttpError
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

func (m *ImageManagerImpl) CreateImage(ctx context.Context, journalExternalId string, image barkat.Image) (*barkat.Image, common.HttpError) {
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Get journal entry to obtain internal ID
		journal, httpErr := m.entryMgr.GetJournal(c, journalExternalId)
		if httpErr != nil {
			return httpErr
		}

		// Set internal ID for foreign key
		image.JournalID = journal.ID

		return m.repo.Create(c, &image)
	})
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (m *ImageManagerImpl) ListImages(ctx context.Context, journalExternalId string) ([]barkat.Image, common.HttpError) {
	var images []barkat.Image
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Get journal entry to obtain internal ID
		journal, httpErr := m.entryMgr.GetJournal(c, journalExternalId)
		if httpErr != nil {
			return httpErr
		}

		// Use internal ID for repository query
		var repoErr common.HttpError
		images, repoErr = m.repo.ListImages(c, journal.ID)
		return repoErr
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (m *ImageManagerImpl) DeleteImage(ctx context.Context, journalExternalId, imageExternalId string) common.HttpError {
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Get journal entry to obtain internal ID
		journal, httpErr := m.entryMgr.GetJournal(c, journalExternalId)
		if httpErr != nil {
			return httpErr
		}
		return m.repo.DeleteById(c, imageExternalId, &barkat.Image{JournalID: journal.ID})
	})
}
