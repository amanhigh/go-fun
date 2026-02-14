package manager

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// TagManager provides business logic for journal tag operations.
type TagManager interface {
	// CreateTag attaches a new tag to an entry.
	CreateTag(ctx context.Context, entryID string, tag *barkat.Tag) common.HttpError
	// ListTags returns all tags for an entry, optionally filtered by type.
	ListTags(ctx context.Context, entryID string, tagType string) ([]barkat.Tag, common.HttpError)
	// DeleteTag removes a tag by ID scoped to an entry.
	DeleteTag(ctx context.Context, entryID string, tagID string) common.HttpError
}

type TagManagerImpl struct {
	entryMgr JournalManager
	repo     repository.TagRepository
}

var _ TagManager = (*TagManagerImpl)(nil)

// NewTagManager creates a new TagManager.
func NewTagManager(entryMgr JournalManager, repo repository.TagRepository) *TagManagerImpl {
	return &TagManagerImpl{entryMgr: entryMgr, repo: repo}
}

func (m *TagManagerImpl) CreateTag(ctx context.Context, entryID string, tag *barkat.Tag) common.HttpError {
	if httpErr := m.entryMgr.EntryExists(ctx, entryID); httpErr != nil {
		return httpErr
	}
	tag.EntryID = entryID
	if err := m.repo.CreateTag(ctx, tag); err != nil {
		return common.NewServerError(fmt.Errorf("failed to create tag: %w", err))
	}
	return nil
}

func (m *TagManagerImpl) ListTags(ctx context.Context, entryID string, tagType string) ([]barkat.Tag, common.HttpError) {
	if httpErr := m.entryMgr.EntryExists(ctx, entryID); httpErr != nil {
		return nil, httpErr
	}
	tags, err := m.repo.ListTags(ctx, entryID, tagType)
	if err != nil {
		return nil, common.NewServerError(fmt.Errorf("failed to list tags: %w", err))
	}
	return tags, nil
}

func (m *TagManagerImpl) DeleteTag(ctx context.Context, entryID string, tagID string) common.HttpError {
	if err := m.repo.DeleteTag(ctx, entryID, tagID); err != nil {
		return common.ErrNotFound
	}
	return nil
}
