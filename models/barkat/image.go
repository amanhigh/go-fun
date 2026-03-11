package barkat

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Image represents a screenshot attached to a journal entry.
type Image struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ExternalID string    `gorm:"column:external_id;uniqueIndex;not null" json:"external_id"`
	JournalID  uint64    `gorm:"column:journal_id;not null;index:idx_image_journal_id" json:"journal_id"`
	Timeframe  string    `gorm:"column:timeframe;not null" json:"timeframe" binding:"required,oneof=DL WK MN TMN SMN YR"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (i *Image) BeforeCreate(_ *gorm.DB) error {
	if i.ExternalID == "" {
		i.ExternalID = "img_" + uuid.NewString()[:8] // Generate external_id with prefix
	}
	if i.CreatedAt.IsZero() {
		i.CreatedAt = time.Now()
	}
	return nil
}
