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
	"github.com/gin-gonic/gin"
)

// NoteHandler provides HTTP handlers for journal note operations.
type NoteHandler interface {
	// HandleCreateNote handles POST /v1/journal-entries/:id/notes
	HandleCreateNote(c *gin.Context)
	// HandleListNotes handles GET /v1/journal-entries/:id/notes
	HandleListNotes(c *gin.Context)
	// HandleDeleteNote handles DELETE /v1/journal-entries/:id/notes/:noteId
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
	entryID := c.Param("id")
	var note barkat.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if httpErr := h.noteMgr.CreateNote(c.Request.Context(), entryID, &note); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, note)
}

func (h *NoteHandlerImpl) HandleListNotes(c *gin.Context) {
	entryID := c.Param("id")
	status := c.Query("note_status")
	notes, httpErr := h.noteMgr.ListNotes(c.Request.Context(), entryID, status)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func (h *NoteHandlerImpl) HandleDeleteNote(c *gin.Context) {
	entryID := c.Param("id")
	noteID := c.Param("noteId")
	if httpErr := h.noteMgr.DeleteNote(c.Request.Context(), entryID, noteID); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.Status(http.StatusNoContent)
}
