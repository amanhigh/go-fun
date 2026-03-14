//nolint:dupl // Intentional CRUD pattern shared with TagHandler for different sub-resources
package handler

// NoteHandler provides HTTP handlers for journal note operations.
// Notes capture trade status snapshots with markdown content.
// This file contains handlers for creating, listing, and deleting journal notes.

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
)

// NoteHandler provides HTTP handlers for journal note operations.
type NoteHandler interface {
	HandleCreateNote(c *gin.Context)
	HandleListNotes(c *gin.Context)
	HandleDeleteNote(c *gin.Context)
}

type NoteHandlerImpl struct {
	noteMgr manager.NoteManager
}

var _ NoteHandler = (*NoteHandlerImpl)(nil)

// NewNoteHandler creates a new NoteHandler.
func NewNoteHandler(noteMgr manager.NoteManager) *NoteHandlerImpl {
	return &NoteHandlerImpl{noteMgr: noteMgr}
}

func (h *NoteHandlerImpl) HandleCreateNote(c *gin.Context) {
	journalID := c.Param("id")
	var note barkat.Note
	if bindErr := c.ShouldBindJSON(&note); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	createdNote, httpErr := h.noteMgr.CreateNote(c.Request.Context(), journalID, note)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, common.NewEnvelope(createdNote))
}

func (h *NoteHandlerImpl) HandleListNotes(c *gin.Context) {
	journalID := c.Param("id")
	status := c.Query("note_status")
	noteList, httpErr := h.noteMgr.ListNotes(c.Request.Context(), journalID, status)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(noteList))
}

func (h *NoteHandlerImpl) HandleDeleteNote(c *gin.Context) {
	journalID := c.Param("id")
	noteID := c.Param("noteId")
	if httpErr := h.noteMgr.DeleteNote(c.Request.Context(), journalID, noteID); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.Status(http.StatusNoContent)
}
