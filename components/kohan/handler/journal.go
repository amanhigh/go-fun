package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/gin-gonic/gin"
)

// JournalHandler provides HTTP handlers for journal entry operations.
//
//go:generate mockery --name JournalHandler
type JournalHandler interface {
	// HandleListEntries handles GET /api/v1/journal-entries
	HandleListEntries(c *gin.Context)
	// HandleGetEntry handles GET /api/v1/journal-entries/:id
	HandleGetEntry(c *gin.Context)
	// HandleCreateEntry handles POST /api/v1/journal-entries
	HandleCreateEntry(c *gin.Context)
}

type JournalHandlerImpl struct {
	journalManager manager.JournalManager
}

var _ JournalHandler = (*JournalHandlerImpl)(nil)

// NewJournalHandler creates a new JournalHandler.
func NewJournalHandler(journalManager manager.JournalManager) *JournalHandlerImpl {
	return &JournalHandlerImpl{journalManager: journalManager}
}

// HandleListEntries handles GET /api/v1/journal-entries
func (h *JournalHandlerImpl) HandleListEntries(c *gin.Context) {
	var query barkat.EntryQuery
	query.Limit = 10 // Default limit

	if err := c.ShouldBindQuery(&query); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	entryList, httpErr := h.journalManager.ListEntries(c.Request.Context(), query)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, entryList)
}

// HandleGetEntry handles GET /api/v1/journal-entries/:id
func (h *JournalHandlerImpl) HandleGetEntry(c *gin.Context) {
	var path barkat.EntryPath
	if err := c.ShouldBindUri(&path); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	entry, httpErr := h.journalManager.GetEntry(c.Request.Context(), path.ID)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, entry)
}

// HandleCreateEntry handles POST /api/v1/journal-entries
func (h *JournalHandlerImpl) HandleCreateEntry(c *gin.Context) {
	var entry barkat.Entry
	if err := c.ShouldBindJSON(&entry); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if httpErr := h.journalManager.CreateEntry(c.Request.Context(), &entry); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, entry)
}
