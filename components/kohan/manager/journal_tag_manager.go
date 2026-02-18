//nolint:dupl
package manager

// TagManager provides business logic for journal tag operations.
// Tags represent categorized labels (reason/management) attached to entries.

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

type TagManager interface {
	// CreateTag attaches a new tag to an entry.
	CreateTag(ctx context.Context, journalID string, tag barkat.Tag) (*barkat.Tag, common.HttpError)
	// ListTags returns all tags for an entry, optionally filtered by type.
	ListTags(ctx context.Context, journalID string, tagType string) ([]barkat.Tag, common.HttpError)
	// DeleteTag removes a tag by ID scoped to an entry.
	DeleteTag(ctx context.Context, journalID string, tagID string) common.HttpError
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

func (m *TagManagerImpl) CreateTag(ctx context.Context, journalID string, tag barkat.Tag) (*barkat.Tag, common.HttpError) {
	if httpErr := m.entryMgr.JournalExists(ctx, journalID); httpErr != nil {
		return nil, httpErr
	}
	tag.JournalID = journalID
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		return m.repo.Create(c, &tag)
	})
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (m *TagManagerImpl) ListTags(ctx context.Context, journalID, tagType string) ([]barkat.Tag, common.HttpError) {
	if httpErr := m.entryMgr.JournalExists(ctx, journalID); httpErr != nil {
		return nil, httpErr
	}

	var tags []barkat.Tag
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		var httpErr common.HttpError
		tags, httpErr = m.repo.ListTags(c, journalID, tagType)
		return httpErr
	})
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (m *TagManagerImpl) DeleteTag(ctx context.Context, journalID, tagID string) common.HttpError {
	if httpErr := m.entryMgr.JournalExists(ctx, journalID); httpErr != nil {
		return httpErr
	}
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		return m.repo.DeleteById(c, tagID, &barkat.Tag{JournalID: journalID})
	})
}
