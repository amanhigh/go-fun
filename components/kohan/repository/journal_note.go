package repository

// NoteRepository provides persistence operations for journal notes.
// Notes capture trade observations and plans at specific entry statuses.

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"gorm.io/gorm"
)

// NoteRepository provides persistence operations for journal notes.
type NoteRepository interface {
	// CreateNote attaches a new note to an entry.
	CreateNote(ctx context.Context, note *barkat.Note) error
	// ListNotes returns all notes for an entry, optionally filtered by status.
	ListNotes(ctx context.Context, entryID string, status string) ([]barkat.Note, error)
	// DeleteNote removes a note by ID scoped to an entry.
	DeleteNote(ctx context.Context, entryID string, noteID string) error
}

type NoteRepositoryImpl struct {
	util.BaseDbRepository
}

var _ NoteRepository = (*NoteRepositoryImpl)(nil)

// NewNoteRepository creates a new NoteRepository backed by GORM.
func NewNoteRepository(db *gorm.DB) *NoteRepositoryImpl {
	return &NoteRepositoryImpl{BaseDbRepository: util.NewBaseDbRepository(db)}
}

func (r *NoteRepositoryImpl) CreateNote(ctx context.Context, note *barkat.Note) error {
	query := r.Db.WithContext(ctx)
	if tx := util.Tx(ctx); tx != nil {
		query = tx
	}
	return util.GormErrorMapper(query.Create(note).Error)
}

func (r *NoteRepositoryImpl) ListNotes(ctx context.Context, entryID, status string) ([]barkat.Note, error) {
	var notes []barkat.Note
	query := r.Db.WithContext(ctx).Where("entry_id = ?", entryID)
	if tx := util.Tx(ctx); tx != nil {
		query = tx
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at").Find(&notes).Error
	return notes, util.GormErrorMapper(err)
}

func (r *NoteRepositoryImpl) DeleteNote(ctx context.Context, entryID, noteID string) error {
	// HACK: Don't create new methods where base can handle it. Remove Simple Methods that don't require ovrride.
	query := r.Db.WithContext(ctx).Where("id = ? AND entry_id = ?", noteID, entryID)
	if tx := util.Tx(ctx); tx != nil {
		query = tx
	}
	result := query.Delete(&barkat.Note{})
	if result.Error != nil {
		return util.GormErrorMapper(result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
