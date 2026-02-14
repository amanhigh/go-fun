package manager

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// JournalManager provides business logic for journal operations.
//
//go:generate mockery --name JournalManager
type JournalManager interface {
	// CreateEntry creates a new journal entry with associations.
	CreateEntry(ctx context.Context, entry *barkat.Entry) common.HttpError
	// GetEntry retrieves a single entry by ID with all associations.
	GetEntry(ctx context.Context, id string) (barkat.Entry, common.HttpError)
	// ListEntries returns a filtered, paginated list of entry summaries.
	ListEntries(ctx context.Context, query barkat.EntryQuery) (barkat.EntryList, common.HttpError)

	// CreateImage attaches a new image to an entry.
	CreateImage(ctx context.Context, entryID string, image *barkat.Image) common.HttpError
	// ListImages returns all images for an entry.
	ListImages(ctx context.Context, entryID string) ([]barkat.Image, common.HttpError)
	// DeleteImage removes an image by ID scoped to an entry.
	DeleteImage(ctx context.Context, entryID string, imageID string) common.HttpError

	// CreateNote attaches a new note to an entry.
	CreateNote(ctx context.Context, entryID string, note *barkat.Note) common.HttpError
	// ListNotes returns all notes for an entry, optionally filtered by status.
	ListNotes(ctx context.Context, entryID string, status string) ([]barkat.Note, common.HttpError)
	// DeleteNote removes a note by ID scoped to an entry.
	DeleteNote(ctx context.Context, entryID string, noteID string) common.HttpError

	// CreateTag attaches a new tag to an entry.
	CreateTag(ctx context.Context, entryID string, tag *barkat.Tag) common.HttpError
	// ListTags returns all tags for an entry, optionally filtered by type.
	ListTags(ctx context.Context, entryID string, tagType string) ([]barkat.Tag, common.HttpError)
	// DeleteTag removes a tag by ID scoped to an entry.
	DeleteTag(ctx context.Context, entryID string, tagID string) common.HttpError
}

type JournalManagerImpl struct {
	repo repository.JournalRepository
}

var _ JournalManager = (*JournalManagerImpl)(nil)

// NewJournalManager creates a new JournalManager.
func NewJournalManager(repo repository.JournalRepository) *JournalManagerImpl {
	return &JournalManagerImpl{repo: repo}
}

// ---- Entry ----

func (m *JournalManagerImpl) CreateEntry(ctx context.Context, entry *barkat.Entry) common.HttpError {
	if err := m.repo.CreateEntry(ctx, entry); err != nil {
		return common.NewServerError(fmt.Errorf("failed to create entry: %w", err))
	}
	return nil
}

func (m *JournalManagerImpl) GetEntry(ctx context.Context, id string) (barkat.Entry, common.HttpError) {
	entry, err := m.repo.GetEntry(ctx, id)
	if err != nil {
		return barkat.Entry{}, common.ErrNotFound
	}
	return entry, nil
}

func (m *JournalManagerImpl) ListEntries(ctx context.Context, query barkat.EntryQuery) (barkat.EntryList, common.HttpError) {
	entries, total, err := m.repo.ListEntries(ctx, query)
	if err != nil {
		return barkat.EntryList{}, common.NewServerError(fmt.Errorf("failed to list entries: %w", err))
	}
	return barkat.EntryList{
		Records:  entries,
		Metadata: common.PaginatedResponse{Total: total},
	}, nil
}

func (m *JournalManagerImpl) checkEntryExists(ctx context.Context, entryID string) common.HttpError {
	exists, err := m.repo.EntryExists(ctx, entryID)
	if err != nil {
		return common.NewServerError(fmt.Errorf("failed to check entry: %w", err))
	}
	if !exists {
		return common.ErrNotFound
	}
	return nil
}

// ---- Image ----

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

// ---- Note ----

func (m *JournalManagerImpl) CreateNote(ctx context.Context, entryID string, note *barkat.Note) common.HttpError {
	if httpErr := m.checkEntryExists(ctx, entryID); httpErr != nil {
		return httpErr
	}
	note.EntryID = entryID
	if err := m.repo.CreateNote(ctx, note); err != nil {
		return common.NewServerError(fmt.Errorf("failed to create note: %w", err))
	}
	return nil
}

func (m *JournalManagerImpl) ListNotes(ctx context.Context, entryID string, status string) ([]barkat.Note, common.HttpError) {
	if httpErr := m.checkEntryExists(ctx, entryID); httpErr != nil {
		return nil, httpErr
	}
	notes, err := m.repo.ListNotes(ctx, entryID, status)
	if err != nil {
		return nil, common.NewServerError(fmt.Errorf("failed to list notes: %w", err))
	}
	return notes, nil
}

func (m *JournalManagerImpl) DeleteNote(ctx context.Context, entryID string, noteID string) common.HttpError {
	if err := m.repo.DeleteNote(ctx, entryID, noteID); err != nil {
		return common.ErrNotFound
	}
	return nil
}

// ---- Tag ----

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
