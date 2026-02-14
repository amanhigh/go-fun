package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/gin-gonic/gin"
)

func (h *JournalHandlerImpl) HandleCreateNote(c *gin.Context) {
	entryID := c.Param("id")
	var note barkat.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if httpErr := h.journalManager.CreateNote(c.Request.Context(), entryID, &note); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, note)
}

func (h *JournalHandlerImpl) HandleListNotes(c *gin.Context) {
	entryID := c.Param("id")
	status := c.Query("note_status")
	notes, httpErr := h.journalManager.ListNotes(c.Request.Context(), entryID, status)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func (h *JournalHandlerImpl) HandleDeleteNote(c *gin.Context) {
	entryID := c.Param("id")
	noteID := c.Param("noteId")
	if httpErr := h.journalManager.DeleteNote(c.Request.Context(), entryID, noteID); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.Status(http.StatusNoContent)
}
