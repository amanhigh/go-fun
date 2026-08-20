//nolint:dupl
package repository

// NoteRepository provides persistence operations for journal notes.
// Notes capture trade observations and plans at specific journal statuses.

import (
	"context"
	"errors"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"gorm.io/gorm"
)

// NoteRepository provides persistence operations for journal notes.
type NoteRepository interface {
	util.BaseDbRepository
	// ListNotes returns all notes for a journal, optionally filtered by status.
	ListNotes(ctx context.Context, journalID uint64, status string) ([]barkat.Note, common.HttpError)
}

type NoteRepositoryImpl struct {
	util.BaseDbRepository
}

var _ NoteRepository = (*NoteRepositoryImpl)(nil)

// NewNoteRepository creates a new NoteRepository backed by GORM.
func NewNoteRepository(baseRepository util.BaseDbRepository) *NoteRepositoryImpl {
	return &NoteRepositoryImpl{BaseDbRepository: baseRepository}
}

func (r *NoteRepositoryImpl) ListNotes(ctx context.Context, journalID uint64, status string) ([]barkat.Note, common.HttpError) {
	var notes []barkat.Note
	var txErr error
	where := barkat.Note{JournalID: journalID}
	if status != "" {
		where.Status = status
	}
	query := r.SafeTx(ctx).Where(&where)
	query = util.ApplySort(query, util.SortOptions{
		DefaultSortBy:    "created_at",
		DefaultSortOrder: common.SortOrderAsc,
	})
	if txErr = query.Find(&notes).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		return nil, util.GormErrorMapper(txErr)
	}
	return notes, nil
}
