package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
)

// JournalHandler provides HTTP handlers for journal entry operations.
type JournalHandler interface {
	HandleListEntries(c *gin.Context)
	HandleGetEntry(c *gin.Context)
	HandleCreateEntry(c *gin.Context)
	HandleDeleteEntry(c *gin.Context)
	HandleUpdateReviewStatus(c *gin.Context)
}

type JournalHandlerImpl struct {
	journalManager manager.JournalManager
}

var _ JournalHandler = (*JournalHandlerImpl)(nil)

// NewJournalHandler creates a new JournalHandler.
func NewJournalHandler(journalManager manager.JournalManager) *JournalHandlerImpl {
	return &JournalHandlerImpl{journalManager: journalManager}
}

// ---- Entry Handlers ----
// FIXME: #A Match other Handlers after Review Comments & Test to Standardize Template.
func (h *JournalHandlerImpl) HandleListEntries(c *gin.Context) {
	query := barkat.NewJournalQuery()

	if bindErr := c.ShouldBindQuery(&query); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	journalList, httpErr := h.journalManager.ListJournals(c.Request.Context(), query)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(journalList))
}

func (h *JournalHandlerImpl) HandleGetEntry(c *gin.Context) {
	var path barkat.JournalPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	journal, httpErr := h.journalManager.GetJournal(c.Request.Context(), path.ID)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(journal))
}

func (h *JournalHandlerImpl) HandleCreateEntry(c *gin.Context) {
	var journal barkat.Journal
	if bindErr := c.ShouldBindJSON(&journal); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	httpErr := h.journalManager.CreateJournal(c.Request.Context(), &journal)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, common.NewEnvelope(&journal))
}

func (h *JournalHandlerImpl) HandleDeleteEntry(c *gin.Context) {
	var path barkat.JournalPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	httpErr := h.journalManager.DeleteJournal(c.Request.Context(), path.ID)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *JournalHandlerImpl) HandleUpdateReviewStatus(c *gin.Context) {
	var path barkat.JournalPath
	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	var update barkat.JournalReviewUpdate
	if bindErr := c.ShouldBindJSON(&update); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	journal, httpErr := h.journalManager.UpdateReviewStatus(c.Request.Context(), path.ID, update)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	// Return minimal response as per PATCH best practices
	response := barkat.UpdateJournalStatusResponse{
		ID:         journal.ExternalID,
		ReviewedAt: journal.ReviewedAt,
	}
	c.JSON(http.StatusOK, common.NewEnvelope(response))
}
