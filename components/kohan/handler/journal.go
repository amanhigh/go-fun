package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/gin-gonic/gin"
)

// JournalHandler provides HTTP handlers for journal entry and sub-resource operations.
//
//go:generate mockery --name JournalHandler
type JournalHandler interface {
	// HandleListEntries handles GET /v1/journal-entries
	HandleListEntries(c *gin.Context)
	// HandleGetEntry handles GET /v1/journal-entries/:id
	HandleGetEntry(c *gin.Context)
	// HandleCreateEntry handles POST /v1/journal-entries
	HandleCreateEntry(c *gin.Context)

	// FIXME: Segregate Image,Notes,Tag into separate handlers, manager,repo for each entity.
	// HandleCreateImage handles POST /v1/journal-entries/:id/images
	HandleCreateImage(c *gin.Context)
	// HandleListImages handles GET /v1/journal-entries/:id/images
	HandleListImages(c *gin.Context)
	// HandleDeleteImage handles DELETE /v1/journal-entries/:id/images/:imageId
	HandleDeleteImage(c *gin.Context)

	// HandleCreateNote handles POST /v1/journal-entries/:id/notes
	HandleCreateNote(c *gin.Context)
	// HandleListNotes handles GET /v1/journal-entries/:id/notes
	HandleListNotes(c *gin.Context)
	// HandleDeleteNote handles DELETE /v1/journal-entries/:id/notes/:noteId
	HandleDeleteNote(c *gin.Context)

	// HandleCreateTag handles POST /v1/journal-entries/:id/tags
	HandleCreateTag(c *gin.Context)
	// HandleListTags handles GET /v1/journal-entries/:id/tags
	HandleListTags(c *gin.Context)
	// HandleDeleteTag handles DELETE /v1/journal-entries/:id/tags/:tagId
	HandleDeleteTag(c *gin.Context)
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

func (h *JournalHandlerImpl) HandleListEntries(c *gin.Context) {
	var query barkat.EntryQuery
	query.Limit = 20

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

// ---- Image Handlers ----

func (h *JournalHandlerImpl) HandleCreateImage(c *gin.Context) {
	entryID := c.Param("id")
	var image barkat.Image
	if err := c.ShouldBindJSON(&image); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if httpErr := h.journalManager.CreateImage(c.Request.Context(), entryID, &image); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, image)
}

func (h *JournalHandlerImpl) HandleListImages(c *gin.Context) {
	entryID := c.Param("id")
	images, httpErr := h.journalManager.ListImages(c.Request.Context(), entryID)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"images": images})
}

func (h *JournalHandlerImpl) HandleDeleteImage(c *gin.Context) {
	entryID := c.Param("id")
	imageID := c.Param("imageId")
	if httpErr := h.journalManager.DeleteImage(c.Request.Context(), entryID, imageID); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.Status(http.StatusNoContent)
}

// ---- Note Handlers ----

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

// ---- Tag Handlers ----

func (h *JournalHandlerImpl) HandleCreateTag(c *gin.Context) {
	entryID := c.Param("id")
	var tag barkat.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if httpErr := h.journalManager.CreateTag(c.Request.Context(), entryID, &tag); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, tag)
}

func (h *JournalHandlerImpl) HandleListTags(c *gin.Context) {
	entryID := c.Param("id")
	tagType := c.Query("type")
	tags, httpErr := h.journalManager.ListTags(c.Request.Context(), entryID, tagType)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

func (h *JournalHandlerImpl) HandleDeleteTag(c *gin.Context) {
	entryID := c.Param("id")
	tagID := c.Param("tagId")
	if httpErr := h.journalManager.DeleteTag(c.Request.Context(), entryID, tagID); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.Status(http.StatusNoContent)
}
