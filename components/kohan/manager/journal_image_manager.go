package manager

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

func (m *JournalManagerImpl) CreateImage(ctx context.Context, entryID string, image *barkat.Image) common.HttpError {
	if httpErr := m.checkEntryExists(ctx, entryID); httpErr != nil {
		return httpErr
	}
	image.EntryID = entryID
	if err := m.repo.CreateImage(ctx, image); err != nil {
		return common.NewServerError(fmt.Errorf("failed to create image: %w", err))
	}
	return nil
}

func (m *JournalManagerImpl) ListImages(ctx context.Context, entryID string) ([]barkat.Image, common.HttpError) {
	if httpErr := m.checkEntryExists(ctx, entryID); httpErr != nil {
		return nil, httpErr
	}
	images, err := m.repo.ListImages(ctx, entryID)
	if err != nil {
		return nil, common.NewServerError(fmt.Errorf("failed to list images: %w", err))
	}
	return images, nil
}

func (m *JournalManagerImpl) DeleteImage(ctx context.Context, entryID string, imageID string) common.HttpError {
	if err := m.repo.DeleteImage(ctx, entryID, imageID); err != nil {
		return common.ErrNotFound
	}
	return nil
}
