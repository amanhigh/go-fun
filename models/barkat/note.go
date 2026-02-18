package barkat

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Note represents a freeform note attached to a journal entry at a specific trade status.
type Note struct {
	ID      string `gorm:"column:id;primaryKey" json:"id"`
	EntryID string `gorm:"column:entry_id;not null;index:idx_note_entry_status,priority:1" json:"entry_id"`
	//nolint:lll // Long validation tag required for oneof
	Status    string    `gorm:"column:status;not null;index:idx_note_entry_status,priority:2" json:"status" binding:"required,oneof=SET RUNNING DROPPED TAKEN REJECTED SUCCESS FAIL MISSED JUST_LOSS BROKEN"`
	Content   string    `gorm:"column:content;not null" json:"content" binding:"required,max=2000"`
	Format    string    `gorm:"column:format;not null;default:markdown" json:"format" binding:"omitempty,oneof=markdown plaintext MARKDOWN PLAINTEXT"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (Note) TableName() string {
	return "journal_notes"
}

func (n *Note) BeforeCreate(_ *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	if n.Format == "" {
		n.Format = "markdown"
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	return nil
}
