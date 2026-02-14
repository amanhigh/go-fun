package manager

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

func (m *JournalManagerImpl) CreateTag(ctx context.Context, entryID string, tag *barkat.Tag) common.HttpError {
	if httpErr := m.checkEntryExists(ctx, entryID); httpErr != nil {
		return httpErr
	}
	tag.EntryID = entryID
	if err := m.repo.CreateTag(ctx, tag); err != nil {
		return common.NewServerError(fmt.Errorf("failed to create tag: %w", err))
	}
	return nil
}

func (m *JournalManagerImpl) ListTags(ctx context.Context, entryID string, tagType string) ([]barkat.Tag, common.HttpError) {
	if httpErr := m.checkEntryExists(ctx, entryID); httpErr != nil {
		return nil, httpErr
	}
	tags, err := m.repo.ListTags(ctx, entryID, tagType)
	if err != nil {
		return nil, common.NewServerError(fmt.Errorf("failed to list tags: %w", err))
	}
	return tags, nil
}

func (m *JournalManagerImpl) DeleteTag(ctx context.Context, entryID string, tagID string) common.HttpError {
	if err := m.repo.DeleteTag(ctx, entryID, tagID); err != nil {
		return common.ErrNotFound
	}
	return nil
}
