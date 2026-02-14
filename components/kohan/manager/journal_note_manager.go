package manager

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

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
