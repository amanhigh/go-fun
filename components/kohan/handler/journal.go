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
// TODO: Match other Handlers after Review Comments & Test to Standardize Template.
func (h *JournalHandlerImpl) HandleListEntries(c *gin.Context) {
	query := barkat.NewJournalQuery()

	if err := c.ShouldBindQuery(&query); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
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

	if err := c.ShouldBindUri(&path); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
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
	if err := c.ShouldBindJSON(&journal); err != nil {
		// FIXME: Enhance ProcessValidationError to return HTTP Code handling Wider Range of Validations.
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
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

	if err := c.ShouldBindUri(&path); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	httpErr := h.journalManager.DeleteJournal(c.Request.Context(), path.ID)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
