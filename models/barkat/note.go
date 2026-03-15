package barkat

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Note represents a freeform note attached to a journal entry at a specific trade status.
type Note struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	ExternalID string    `gorm:"column:external_id;uniqueIndex;not null" json:"id"`
	JournalID  uint64    `gorm:"column:journal_id;not null;index:idx_note_journal_id" json:"journal_id"`
	Status     string    `gorm:"column:status;not null;index:idx_note_status" json:"status" binding:"required,oneof=SET RUNNING DROPPED TAKEN REJECTED SUCCESS FAIL MISSED JUST_LOSS BROKEN"`
	Content    string    `gorm:"column:content;not null" json:"content" binding:"required,max=2000"`
	Format     string    `gorm:"column:format;not null;default:MARKDOWN" json:"format" binding:"omitempty,oneof=MARKDOWN PLAINTEXT"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (n *Note) BeforeCreate(_ *gorm.DB) error {
	if n.ExternalID == "" {
		n.ExternalID = "not_" + uuid.NewString()[:8] // Generate external_id with prefix
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	return nil
}

// NoteList is the response for notes.
type NoteList struct {
	Notes []Note `json:"notes"`
}

// NoteQuery holds query parameters for listing/filtering notes.
type NoteQuery struct {
	Status string `form:"note_status" binding:"omitempty,oneof=SET RUNNING DROPPED TAKEN REJECTED SUCCESS FAIL MISSED JUST_LOSS BROKEN"`
}

// ---- Path Parameter Structs ----

// NotePath binds the :id and :noteId path parameters for note operations.
type NotePath struct {
	JournalID string `uri:"id" binding:"required,journal_id"`
	NoteID    string `uri:"noteId" binding:"required,note_id"`
}
