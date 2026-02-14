package barkat

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tag represents a reason or management tag attached to a journal entry.
type Tag struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	EntryID   string    `gorm:"column:entry_id;not null;index:idx_tag_entry_type,priority:1" json:"entry_id"`
	Tag       string    `gorm:"column:tag;not null;index:idx_tag_type_value,priority:2" json:"tag" binding:"required"`
	Type      string    `gorm:"column:type;not null;index:idx_tag_entry_type,priority:2;index:idx_tag_type_value,priority:1" json:"type" binding:"required,oneof=reason management"`
	Override  *string   `gorm:"column:override" json:"override,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (Tag) TableName() string {
	return "journal_tags"
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
