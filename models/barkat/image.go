package barkat

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Image represents a screenshot attached to a journal entry.
type Image struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	JournalID string    `gorm:"column:journal_id;not null;index:idx_image_journal_timeframe,priority:1" json:"journal_id"`
	Timeframe string    `gorm:"column:timeframe;not null;index:idx_image_journal_timeframe,priority:2" json:"timeframe" binding:"required,oneof=DL WK MN TMN SMN YR"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
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
