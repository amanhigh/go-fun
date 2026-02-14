package barkat

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Image represents a screenshot attached to a journal entry.
type Image struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	EntryID   string    `gorm:"column:entry_id;not null" json:"entry_id"`
	Position  int       `gorm:"column:position;not null" json:"position"`
	Path      string    `gorm:"column:path;not null" json:"path"`
	IsCheck   bool      `gorm:"column:is_check;not null;default:false" json:"is_check"`
	Width     *int      `gorm:"column:width" json:"width,omitempty"`
	Height    *int      `gorm:"column:height" json:"height,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (Image) TableName() string {
	return "barkat_images"
}

func (i *Image) BeforeCreate(_ *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	if i.CreatedAt.IsZero() {
		i.CreatedAt = time.Now()
	}
	return nil
}
