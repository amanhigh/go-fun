package barkat

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Entry represents a single trade journal capture event.
type Entry struct {
	ID            string     `gorm:"column:id;primaryKey" json:"id"`
	Ticker        string     `gorm:"column:ticker;not null" json:"ticker"`
	Sequence      string     `gorm:"column:sequence;not null" json:"sequence"`
	Type          string     `gorm:"column:type;not null" json:"type"`
	Outcome       string     `gorm:"column:outcome;not null" json:"outcome"`
	Trend         string     `gorm:"column:trend;not null;default:trend" json:"trend"`
	NotesMarkdown *string    `gorm:"column:notes_markdown" json:"notes_markdown,omitempty"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Associations
	Images []Image `gorm:"foreignKey:EntryID;references:ID" json:"images,omitempty"`
}

func (Entry) TableName() string {
	return "barkat_entries"
}

func (e *Entry) BeforeCreate(_ *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	return nil
}

// EntryPath binds the :id path parameter.
type EntryPath struct {
	ID string `uri:"id" binding:"required"`
}

// EntryQuery holds query parameters for listing/filtering entries.
type EntryQuery struct {
	common.Pagination
	Ticker   string `form:"ticker" binding:"omitempty,min=1,max=30"`
	Type     string `form:"type" binding:"omitempty,oneof=rejected result set"`
	Outcome  string `form:"outcome" binding:"omitempty,oneof=fail broken taken success running justloss"`
	Sequence string `form:"sequence" binding:"omitempty,oneof=mwd wdh yr"`
	From     string `form:"from" binding:"omitempty"`
	To       string `form:"to" binding:"omitempty"`
}

// EntryList is the paginated response for journal entries.
type EntryList struct {
	Records  []Entry                  `json:"records"`
	Metadata common.PaginatedResponse `json:"metadata"`
}
