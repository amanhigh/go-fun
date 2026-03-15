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
	var path barkat.JournalPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	var note barkat.Note
	if bindErr := c.ShouldBindJSON(&note); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	createdNote, httpErr := h.noteMgr.CreateNote(c.Request.Context(), path.JournalID, note)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, common.NewEnvelope(createdNote))
}

func (h *NoteHandlerImpl) HandleListNotes(c *gin.Context) {
	var path barkat.JournalPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	query := barkat.NoteQuery{}

	if bindErr := c.ShouldBindQuery(&query); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	noteList, httpErr := h.noteMgr.ListNotes(c.Request.Context(), path.JournalID, query.Status)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(noteList))
}

func (h *NoteHandlerImpl) HandleDeleteNote(c *gin.Context) {
	var path barkat.NotePath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	if httpErr := h.noteMgr.DeleteNote(c.Request.Context(), path.JournalID, path.NoteID); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.Status(http.StatusNoContent)
}
