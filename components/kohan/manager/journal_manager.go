package manager

import (
	"context"
	"fmt"
	"net/http"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// Validation constants per PRD specifications
const (
	MinImages      = 4
	MaxImages      = 6
	MaxNotes       = 1
	MaxNoteContent = 2000
	MaxTagLength   = 10
	MaxTagOverride = 5
)

// JournalManager provides business logic for journal entry operations.
type JournalManager interface {
	// CreateEntry creates a new journal entry with associations.
	CreateEntry(ctx context.Context, entry *barkat.Entry) common.HttpError
	// GetEntry retrieves a single entry by ID with all associations.
	GetEntry(ctx context.Context, id string) (barkat.Entry, common.HttpError)
	// ListEntries returns a filtered, paginated list of entry summaries.
	ListEntries(ctx context.Context, query barkat.EntryQuery) (barkat.EntryList, common.HttpError)
	// EntryExists checks if an entry with the given ID exists.
	EntryExists(ctx context.Context, entryID string) common.HttpError
	// DeleteEntry deletes a journal entry by ID.
	DeleteEntry(ctx context.Context, id string) common.HttpError
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

//nolint:gocyclo,cyclop,funlen // Validation logic requires multiple checks per PRD
func (m *JournalManagerImpl) CreateEntry(ctx context.Context, entry *barkat.Entry) common.HttpError {
	// Validate images: PRD requires min 4, max 6
	if len(entry.Images) < MinImages {
		return common.NewHttpError(fmt.Sprintf("images: minimum %d required", MinImages), http.StatusBadRequest)
	}
	if len(entry.Images) > MaxImages {
		return common.NewHttpError(fmt.Sprintf("images: maximum %d allowed", MaxImages), http.StatusRequestEntityTooLarge)
	}

	// Validate each image timeframe
	validTimeframes := map[string]bool{"DL": true, "WK": true, "MN": true, "TMN": true, "SMN": true, "YR": true}
	for _, img := range entry.Images {
		if !validTimeframes[img.Timeframe] {
			return common.NewHttpError("images.timeframe: must be one of DL,WK,MN,TMN,SMN,YR", http.StatusBadRequest)
		}
	}

	// Validate notes: PRD allows max 1 at create time
	if len(entry.Notes) > MaxNotes {
		return common.NewHttpError(fmt.Sprintf("note_blocks: maximum %d allowed at create", MaxNotes), http.StatusRequestEntityTooLarge)
	}

	// Validate each note
	validNoteStatuses := map[string]bool{
		"SET": true, "RUNNING": true, "DROPPED": true, "TAKEN": true, "REJECTED": true,
		"SUCCESS": true, "FAIL": true, "MISSED": true, "JUST_LOSS": true, "BROKEN": true,
	}
	validNoteFormats := map[string]bool{"MARKDOWN": true, "PLAINTEXT": true, "markdown": true, "plaintext": true, "": true}
	for _, note := range entry.Notes {
		if !validNoteStatuses[note.Status] {
			return common.NewHttpError("note_blocks.status: invalid status", http.StatusBadRequest)
		}
		if !validNoteFormats[note.Format] {
			return common.NewHttpError("note_blocks.format: must be MARKDOWN or PLAINTEXT", http.StatusBadRequest)
		}
		if len(note.Content) > MaxNoteContent {
			return common.NewHttpError(fmt.Sprintf("note_blocks.content: maximum %d characters", MaxNoteContent), http.StatusBadRequest)
		}
	}

	// Validate tags
	validTagTypes := map[string]bool{"REASON": true, "MANAGEMENT": true, "reason": true, "management": true}
	for _, tag := range entry.Tags {
		if !validTagTypes[tag.Type] {
			return common.NewHttpError("tags.type: must be REASON or MANAGEMENT", http.StatusBadRequest)
		}
		if len(tag.Tag) > MaxTagLength {
			return common.NewHttpError(fmt.Sprintf("tags.tag: maximum %d characters", MaxTagLength), http.StatusBadRequest)
		}
		if tag.Override != nil && len(*tag.Override) > MaxTagOverride {
			return common.NewHttpError(fmt.Sprintf("tags.override: maximum %d characters", MaxTagOverride), http.StatusBadRequest)
		}
	}

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

func (m *JournalManagerImpl) EntryExists(ctx context.Context, entryID string) common.HttpError {
	entry := &barkat.Entry{}
	return m.repo.FindById(ctx, entryID, entry)
}

func (m *JournalManagerImpl) DeleteEntry(ctx context.Context, id string) common.HttpError {
	entry := &barkat.Entry{}
	return m.repo.DeleteById(ctx, id, entry)
}
