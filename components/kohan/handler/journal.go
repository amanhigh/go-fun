package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
	"github.com/golang-sql/civil"
)

// JournalHandler provides HTTP handlers for journal operations.
type JournalHandler interface {
	HandleListJournals(c *gin.Context)
	HandleGetJournal(c *gin.Context)
	HandleCreateJournal(c *gin.Context)
	HandleDeleteJournal(c *gin.Context)
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

// ---- Journal Handlers ----
func (h *JournalHandlerImpl) HandleListJournals(c *gin.Context) {
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

func (h *JournalHandlerImpl) HandleGetJournal(c *gin.Context) {
	var path barkat.JournalPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	journal, httpErr := h.journalManager.GetJournal(c.Request.Context(), path.JournalID)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(journal))
}

func (h *JournalHandlerImpl) HandleCreateJournal(c *gin.Context) {
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

func (h *JournalHandlerImpl) HandleDeleteJournal(c *gin.Context) {
	var path barkat.JournalPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	if httpErr := h.journalManager.DeleteJournal(c.Request.Context(), path.JournalID); httpErr != nil {
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
	if err := c.ShouldBindJSON(&update); err != nil {
		httpErr := util.ProcessValidationError(err)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	journal, httpErr := h.journalManager.UpdateReviewStatus(c.Request.Context(), path.JournalID, update)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	// Return minimal response with id, status, and reviewed_at as per PRD
	var reviewedAt *civil.Date
	if journal.ReviewedAt != nil {
		civilDate := civil.DateOf(*journal.ReviewedAt)
		reviewedAt = &civilDate
	}
	response := barkat.UpdateJournalStatusResponse{
		ID:         journal.ExternalID,
		Status:     journal.Status,
		ReviewedAt: reviewedAt,
	}
	c.JSON(http.StatusOK, common.NewEnvelope(response))
}
