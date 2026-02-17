package barkat

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Entry represents a single trade journal capture event.
// BUG: Rename to Journal all entities in this file & filename.
type Entry struct {
	ID        string     `gorm:"column:id;primaryKey" json:"id"`
	Ticker    string     `gorm:"column:ticker;not null" json:"ticker" binding:"required,max=10"`
	Sequence  string     `gorm:"column:sequence;not null" json:"sequence" binding:"required,oneof=MWD YR"`
	Type      string     `gorm:"column:type;not null" json:"type" binding:"required,oneof=REJECTED RESULT SET"`
	Status    string     `gorm:"column:status;not null" json:"status" binding:"required,oneof=SET RUNNING DROPPED TAKEN REJECTED SUCCESS FAIL MISSED JUST_LOSS BROKEN"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;index:idx_entry_ticker_created,priority:2,sort:desc" json:"created_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Associations
	Images []Image `gorm:"foreignKey:EntryID;references:ID" json:"images,omitempty" binding:"required,min=4,max=6,dive"`
	Tags   []Tag   `gorm:"foreignKey:EntryID;references:ID" json:"tags,omitempty" binding:"dive"`
	Notes  []Note  `gorm:"foreignKey:EntryID;references:ID" json:"notes,omitempty" binding:"max=1,dive"`
}

func (Entry) TableName() string {
	return "journal"
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
	Ticker        string `form:"ticker" binding:"omitempty,min=1,max=10"`
	Type          string `form:"type" binding:"omitempty,oneof=REJECTED RESULT SET"`
	Status        string `form:"status" binding:"omitempty,oneof=SET RUNNING DROPPED TAKEN REJECTED SUCCESS FAIL MISSED JUST_LOSS BROKEN"`
	Sequence      string `form:"sequence" binding:"omitempty,oneof=MWD YR"`
	CreatedAfter  string `form:"created-after" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	CreatedBefore string `form:"created-before" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	SortBy        string `form:"sort-by" binding:"omitempty,oneof=created_at ticker sequence"`
	SortOrder     string `form:"sort-order" binding:"omitempty,oneof=asc desc"`
}

// EntryList is the paginated response for journal entries.
type EntryList struct {
	Records []Entry `json:"journals"`
	// BUG: Inline Pagination Metadat as per PRD ?
	Metadata common.PaginatedResponse `json:"metadata"`
}
