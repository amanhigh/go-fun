package barkat

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Journal API route constants
const (
	// Base journal routes
	JournalBase      = common.APIV1 + "/journals"
	JournalEntries   = JournalBase
	JournalEntryByID = JournalBase + "/:id"

	// Journal sub-resource routes
	JournalImages    = JournalEntryByID + "/images"
	JournalImageByID = JournalImages + "/:imageId"
	JournalNotes     = JournalEntryByID + "/notes"
	JournalNoteByID  = JournalNotes + "/:noteId"
	JournalTags      = JournalEntryByID + "/tags"
	JournalTagByID   = JournalTags + "/:tagId"
)

// Journal represents a single trade journal capture event.
type Journal struct {
	ID         uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	ExternalID string     `gorm:"column:external_id;uniqueIndex;not null" json:"id"`
	Ticker     string     `gorm:"column:ticker;not null;index:idx_journal_ticker" json:"ticker" binding:"required,max=10,ticker"`
	Sequence   string     `gorm:"column:sequence;not null" json:"sequence" binding:"required,oneof=MWD YR"`
	Type       string     `gorm:"column:type;not null" json:"type" binding:"required,oneof=REJECTED RESULT SET"`
	Status     string     `gorm:"column:status;not null" json:"status" binding:"required,oneof=SET RUNNING DROPPED TAKEN REJECTED SUCCESS FAIL MISSED JUST_LOSS BROKEN"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null;index:idx_journal_created_at" json:"created_at"`
	ReviewedAt *time.Time `gorm:"column:reviewed_at" json:"reviewed_at,omitempty"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Associations
	Images []Image `gorm:"foreignKey:JournalID;references:ID" json:"images,omitempty" binding:"required,min=4,max=6,dive"`
	Tags   []Tag   `gorm:"foreignKey:JournalID;references:ID" json:"tags,omitempty" binding:"max=10,dive"`
	Notes  []Note  `gorm:"foreignKey:JournalID;references:ID" json:"notes,omitempty" binding:"max=1,dive"`
}

func (j *Journal) BeforeCreate(_ *gorm.DB) error {
	if j.ExternalID == "" {
		j.ExternalID = "jrn_" + uuid.NewString()[:8] // Generate external_id with prefix
	}
	return nil
}

// JournalPath binds the :id path parameter.
type JournalPath struct {
	ID string `uri:"id" binding:"required"`
}

// JournalQuery holds query parameters for listing/filtering journals.
type JournalQuery struct {
	common.Pagination
	Ticker        string `form:"ticker" binding:"omitempty,min=1,max=10,ticker"`
	Type          string `form:"type" binding:"omitempty,oneof=REJECTED RESULT SET"`
	Status        string `form:"status" binding:"omitempty,oneof=SET RUNNING DROPPED TAKEN REJECTED SUCCESS FAIL MISSED JUST_LOSS BROKEN"`
	Sequence      string `form:"sequence" binding:"omitempty,oneof=MWD YR"`
	CreatedAfter  string `form:"created-after" binding:"omitempty,datetime=2006-01-02"`
	CreatedBefore string `form:"created-before" binding:"omitempty,datetime=2006-01-02"`
	Reviewed      *bool  `form:"reviewed" binding:"omitempty"`
	SortBy        string `form:"sort-by" binding:"omitempty,oneof=created_at ticker sequence"`
	SortOrder     string `form:"sort-order" binding:"omitempty,oneof=asc desc"`
}

// NewJournalQuery creates a JournalQuery struct with default pagination values
func NewJournalQuery() JournalQuery {
	return JournalQuery{
		Pagination: common.Pagination{},
	}
}

// JournalList is the paginated response for journals.
type JournalList struct {
	Journals []Journal                `json:"journals"`
	Metadata common.PaginatedResponse `json:"metadata"`
}

// JournalReviewUpdate represents the request body for updating journal review status.
type JournalReviewUpdate struct {
	ReviewedAt string `json:"reviewed_at" binding:"required,datetime=2006-01-02,future_date"`
}

// UpdateJournalStatusResponse represents the response for PATCH review status updates.
// This follows the PRD specification for minimal PATCH responses.
type UpdateJournalStatusResponse struct {
	ID         string  `json:"id"`
	ReviewedAt *string `json:"reviewed_at"`
}
