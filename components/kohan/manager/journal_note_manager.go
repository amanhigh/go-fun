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
	ListNotes(ctx context.Context, journalID string, status string) (barkat.NoteList, common.HttpError)
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

func (m *NoteManagerImpl) CreateNote(ctx context.Context, journalExternalId string, note barkat.Note) (*barkat.Note, common.HttpError) {
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Get journal entry to obtain internal ID
		journal, httpErr := m.entryMgr.GetJournal(c, journalExternalId)
		if httpErr != nil {
			return httpErr
		}

		// Set internal ID for foreign key
		note.JournalID = journal.ID

		return m.repo.Create(c, &note)
	})
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (m *NoteManagerImpl) ListNotes(ctx context.Context, journalID, status string) (barkat.NoteList, common.HttpError) {
	var notes []barkat.Note
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Get journal entry to obtain internal ID
		journal, httpErr := m.entryMgr.GetJournal(c, journalID)
		if httpErr != nil {
			return httpErr
		}

		// Use internal ID for repository query
		var repoErr common.HttpError
		notes, repoErr = m.repo.ListNotes(c, journal.ID, status)
		return repoErr
	})
	if err != nil {
		return barkat.NoteList{}, err
	}
	return barkat.NoteList{Notes: notes}, nil
}

func (m *NoteManagerImpl) DeleteNote(ctx context.Context, journalExternalId, noteExternalId string) common.HttpError {
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Get journal entry to obtain internal ID
		journal, httpErr := m.entryMgr.GetJournal(c, journalExternalId)
		if httpErr != nil {
			return httpErr
		}

		// First fetch the note by external_id to get internal ID
		var note barkat.Note
		httpErr = m.repo.GetByExternalId(c, noteExternalId, &note)
		if httpErr != nil {
			return httpErr
		}

		// Verify the note belongs to the correct journal
		if note.JournalID != journal.ID {
			return common.ErrNotFound
		}

		// Now delete by internal ID using base repository method
		return m.repo.DeleteById(c, note.ID, &barkat.Note{})
	})
}
