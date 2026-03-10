package barkat

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tag represents a reason or management tag attached to a journal entry.
type Tag struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	JournalID string    `gorm:"column:journal_id;not null" json:"journal_id"`
	Tag       string    `gorm:"column:tag;not null" json:"tag" binding:"required,max=10,tag_format"`
	Type      string    `gorm:"column:type;not null" json:"type" binding:"required,oneof=REASON MANAGEMENT"`
	Override  *string   `gorm:"column:override" json:"override,omitempty" binding:"omitempty,max=5,override_format"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (t *Tag) BeforeCreate(_ *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	return nil
}
