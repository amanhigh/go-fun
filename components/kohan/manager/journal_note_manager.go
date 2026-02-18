//nolint:dupl
package manager

// NoteManager provides business logic for journal note operations.
// Notes represent freeform text attached to entries at specific trade statuses.

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

type NoteManager interface {
	// CreateNote attaches a new note to an entry.
	CreateNote(ctx context.Context, journalID string, note barkat.Note) (*barkat.Note, common.HttpError)
	// ListNotes returns all notes for an entry, optionally filtered by status.
	ListNotes(ctx context.Context, journalID string, status string) ([]barkat.Note, common.HttpError)
	// DeleteNote removes a note by ID scoped to an entry.
	DeleteNote(ctx context.Context, journalID string, noteID string) common.HttpError
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

func (m *NoteManagerImpl) CreateNote(ctx context.Context, journalID string, note barkat.Note) (*barkat.Note, common.HttpError) {
	if httpErr := m.entryMgr.JournalExists(ctx, journalID); httpErr != nil {
		return nil, httpErr
	}
	note.JournalID = journalID
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		return m.repo.Create(c, &note)
	})
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (m *NoteManagerImpl) ListNotes(ctx context.Context, journalID, status string) ([]barkat.Note, common.HttpError) {
	if httpErr := m.entryMgr.JournalExists(ctx, journalID); httpErr != nil {
		return nil, httpErr
	}

	var notes []barkat.Note
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		var httpErr common.HttpError
		notes, httpErr = m.repo.ListNotes(c, journalID, status)
		return httpErr
	})
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (m *NoteManagerImpl) DeleteNote(ctx context.Context, journalID, noteID string) common.HttpError {
	if httpErr := m.entryMgr.JournalExists(ctx, journalID); httpErr != nil {
		return httpErr
	}
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		return m.repo.DeleteById(c, noteID, &barkat.Note{JournalID: journalID})
	})
}
