package repository

import (
	"context"

	"github.com/amanhigh/go-fun/models/barkat"
	"gorm.io/gorm"
)

func (r *JournalRepositoryImpl) CreateNote(ctx context.Context, note *barkat.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *JournalRepositoryImpl) ListNotes(ctx context.Context, entryID string, status string) ([]barkat.Note, error) {
	var notes []barkat.Note
	tx := r.db.WithContext(ctx).Where("entry_id = ?", entryID)
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	err := tx.Order("created_at").Find(&notes).Error
	return notes, err
}

func (r *JournalRepositoryImpl) DeleteNote(ctx context.Context, entryID string, noteID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND entry_id = ?", noteID, entryID).Delete(&barkat.Note{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
