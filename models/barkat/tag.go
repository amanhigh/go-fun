package barkat

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tag represents a reason or management tag attached to a journal entry.
type Tag struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ExternalID string    `gorm:"column:external_id;uniqueIndex;not null" json:"external_id"`
	JournalID  uint64    `gorm:"column:journal_id;not null;index:idx_tag_journal_id" json:"journal_id"`
	Tag        string    `gorm:"column:tag;not null" json:"tag" binding:"required,max=10,tag"`
	Type       string    `gorm:"column:type;not null" json:"type" binding:"required,oneof=REASON MANAGEMENT"`
	Override   *string   `gorm:"column:override" json:"override,omitempty" binding:"omitempty,max=5,override"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (t *Tag) BeforeCreate(_ *gorm.DB) error {
	if t.ExternalID == "" {
		t.ExternalID = "tag_" + uuid.NewString()[:8] // Generate external_id with prefix
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	return nil
}
