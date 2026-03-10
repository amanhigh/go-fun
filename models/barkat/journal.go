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
	ID        string     `gorm:"column:id;primaryKey" json:"id"`
	Ticker    string     `gorm:"column:ticker;not null" json:"ticker" binding:"required,max=10,ticker"`
	Sequence  string     `gorm:"column:sequence;not null" json:"sequence" binding:"required,oneof=MWD YR"`
	Type      string     `gorm:"column:type;not null" json:"type" binding:"required,oneof=REJECTED RESULT SET"`
	Status    string     `gorm:"column:status;not null" json:"status" binding:"required,oneof=SET RUNNING DROPPED TAKEN REJECTED SUCCESS FAIL MISSED JUST_LOSS BROKEN"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;index:idx_journal_ticker_created,priority:2,sort:desc" json:"created_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Associations
	Images []Image `gorm:"foreignKey:JournalID;references:ID" json:"images,omitempty" binding:"required,min=4,max=6,dive"`
	Tags   []Tag   `gorm:"foreignKey:JournalID;references:ID" json:"tags,omitempty" binding:"max=10,dive"`
	Notes  []Note  `gorm:"foreignKey:JournalID;references:ID" json:"notes,omitempty" binding:"max=1,dive"`
}

func (j *Journal) BeforeCreate(_ *gorm.DB) error {
	if j.ID == "" {
		j.ID = uuid.NewString()
	}
	// Set JournalID for all associated images (required for GORM associations)
	for i := range j.Images {
		// HACK: FK on Internal Id or external id ?
		j.Images[i].JournalID = j.ID
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
	Records  []Journal                `json:"journals"`
	Metadata common.PaginatedResponse `json:"metadata"`
}
