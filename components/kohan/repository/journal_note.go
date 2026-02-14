package repository

import (
	"context"

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
	db *gorm.DB
}

var _ NoteRepository = (*NoteRepositoryImpl)(nil)

// NewNoteRepository creates a new NoteRepository backed by GORM.
func NewNoteRepository(db *gorm.DB) *NoteRepositoryImpl {
	return &NoteRepositoryImpl{db: db}
}

func (r *NoteRepositoryImpl) CreateNote(ctx context.Context, note *barkat.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *NoteRepositoryImpl) ListNotes(ctx context.Context, entryID string, status string) ([]barkat.Note, error) {
	var notes []barkat.Note
	tx := r.db.WithContext(ctx).Where("entry_id = ?", entryID)
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	err := tx.Order("created_at").Find(&notes).Error
	return notes, err
}

func (r *NoteRepositoryImpl) DeleteNote(ctx context.Context, entryID string, noteID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND entry_id = ?", noteID, entryID).Delete(&barkat.Note{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
