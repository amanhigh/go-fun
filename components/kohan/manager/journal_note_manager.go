package manager

// NoteManager provides business logic for journal note operations.
// Notes represent freeform text attached to entries at specific trade statuses.

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// NoteManager provides business logic for journal note operations.
type NoteManager interface {
	// CreateNote attaches a new note to an entry.
	CreateNote(ctx context.Context, entryID string, note *barkat.Note) common.HttpError
	// ListNotes returns all notes for an entry, optionally filtered by status.
	ListNotes(ctx context.Context, entryID string, status string) ([]barkat.Note, common.HttpError)
	// DeleteNote removes a note by ID scoped to an entry.
	DeleteNote(ctx context.Context, entryID string, noteID string) common.HttpError
}

type NoteManagerImpl struct {
	entryMgr JournalManager
	repo     repository.NoteRepository
}

var _ NoteManager = (*NoteManagerImpl)(nil)

// NewNoteManager creates a new NoteManager.
func NewNoteManager(entryMgr JournalManager, repo repository.NoteRepository) *NoteManagerImpl {
	return &NoteManagerImpl{entryMgr: entryMgr, repo: repo}
}

func (m *NoteManagerImpl) CreateNote(ctx context.Context, entryID string, note *barkat.Note) common.HttpError {
	if httpErr := m.entryMgr.EntryExists(ctx, entryID); httpErr != nil {
		return httpErr
	}
	note.EntryID = entryID

	// Handle transactions at manager layer using the embedded BaseDbRepository
	noteRepoImpl, ok := m.repo.(*repository.NoteRepositoryImpl)
	if !ok {
		return common.NewServerError(fmt.Errorf("failed to cast repository to NoteRepositoryImpl"))
	}
	return noteRepoImpl.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		if err := m.repo.CreateNote(c, note); err != nil {
			return common.NewServerError(fmt.Errorf("failed to create note: %w", err))
		}
		return nil
	})
}

func (m *NoteManagerImpl) ListNotes(ctx context.Context, entryID, status string) ([]barkat.Note, common.HttpError) {
	if httpErr := m.entryMgr.EntryExists(ctx, entryID); httpErr != nil {
		return nil, httpErr
	}
	notes, err := m.repo.ListNotes(ctx, entryID, status)
	if err != nil {
		return nil, common.NewServerError(fmt.Errorf("failed to list notes: %w", err))
	}
	return notes, nil
}

func (m *NoteManagerImpl) DeleteNote(ctx context.Context, entryID, noteID string) common.HttpError {
	if err := m.repo.DeleteNote(ctx, entryID, noteID); err != nil {
		return common.ErrNotFound
	}
	return nil
}
