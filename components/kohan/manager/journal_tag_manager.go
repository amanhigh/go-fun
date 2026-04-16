//nolint:dupl
package manager

// TagManager provides business logic for journal tag operations.
// Tags represent categorized labels (reason/management) attached to entries.

// HACK: No DB backup strategy for journal SQLite file at data/journals.db
// Consider adding automated backup to object storage (S3-compatible) with retention

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

type TagManager interface {
	// CreateTag attaches a new tag to a journal.
	CreateTag(ctx context.Context, journalID string, tag barkat.Tag) (*barkat.Tag, common.HttpError)
	// ListTags returns all tags for a journal, optionally filtered by type.
	ListTags(ctx context.Context, journalID string, tagType string) (barkat.TagList, common.HttpError)
	// DeleteTag removes a tag by ID scoped to a journal.
	DeleteTag(ctx context.Context, journalID string, tagID string) common.HttpError
}

type TagManagerImpl struct {
	journalMgr JournalManager
	repo       repository.TagRepository
}

var _ TagManager = (*TagManagerImpl)(nil)

// NewTagManager creates a new TagManager.
func NewTagManager(journalMgr JournalManager, repo repository.TagRepository) *TagManagerImpl {
	return &TagManagerImpl{journalMgr: journalMgr, repo: repo}
}

func (m *TagManagerImpl) CreateTag(ctx context.Context, journalExternalId string, tag barkat.Tag) (*barkat.Tag, common.HttpError) {
	// FIXME: Add explicit allowed-values validation for tags and overrides before persisting
	// - Validate tag value is in the allowed set (e.g., "dep", "nca", "oe", "ntr", "important")
	// - Validate override value is in the allowed set (e.g., "loc", "egf", "loc1", "egf1")
	// - Validate override is only provided for REASON type tags (not MANAGEMENT/DIRECTION)
	// Currently relies on format validators (tagRegex/overrideRegex) only, not whitelist validation

	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Get journal to obtain internal ID
		journal, httpErr := m.journalMgr.GetJournal(c, journalExternalId)
		if httpErr != nil {
			return httpErr
		}

		// Set internal ID for foreign key
		tag.JournalID = journal.ID

		return m.repo.Create(c, &tag)
	})
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (m *TagManagerImpl) ListTags(ctx context.Context, journalID, tagType string) (barkat.TagList, common.HttpError) {
	var tags []barkat.Tag
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Get journal to obtain internal ID
		journal, httpErr := m.journalMgr.GetJournal(c, journalID)
		if httpErr != nil {
			return httpErr
		}

		// Use internal ID for repository query
		var repoErr common.HttpError
		tags, repoErr = m.repo.ListTags(c, journal.ID, tagType)
		return repoErr
	})
	if err != nil {
		return barkat.TagList{}, err
	}
	return barkat.TagList{Tags: tags}, nil
}

func (m *TagManagerImpl) DeleteTag(ctx context.Context, journalExternalId, tagExternalId string) common.HttpError {
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Get journal to obtain internal ID
		journal, httpErr := m.journalMgr.GetJournal(c, journalExternalId)
		if httpErr != nil {
			return httpErr
		}

		// First fetch the tag by external_id to get internal ID
		var tag barkat.Tag
		httpErr = m.repo.GetByExternalId(c, tagExternalId, &tag)
		if httpErr != nil {
			return httpErr
		}

		// Verify the tag belongs to the correct journal
		if tag.JournalID != journal.ID {
			return common.ErrNotFound
		}

		// Now delete by internal ID using base repository method
		return m.repo.DeleteById(c, tag.ID, &barkat.Tag{})
	})
}
