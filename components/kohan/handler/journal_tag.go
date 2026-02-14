//nolint:dupl // Intentional CRUD pattern shared with NoteHandler for different sub-resources
package handler

// TagHandler provides HTTP handlers for journal tag operations.
// Tags categorize entries with reason codes (oe, dep, etc.) or management labels.

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/gin-gonic/gin"
)

// TagHandler provides HTTP handlers for journal tag operations.
type TagHandler interface {
	// HandleCreateTag handles POST /v1/journal-entries/:id/tags
	HandleCreateTag(c *gin.Context)
	// HandleListTags handles GET /v1/journal-entries/:id/tags
	HandleListTags(c *gin.Context)
	// HandleDeleteTag handles DELETE /v1/journal-entries/:id/tags/:tagId
	HandleDeleteTag(c *gin.Context)
}

type TagHandlerImpl struct {
	tagMgr manager.TagManager
}

var _ TagHandler = (*TagHandlerImpl)(nil)

// NewTagHandler creates a new TagHandler.
func NewTagHandler(tagMgr manager.TagManager) *TagHandlerImpl {
	return &TagHandlerImpl{tagMgr: tagMgr}
}

func (h *TagHandlerImpl) HandleCreateTag(c *gin.Context) {
	entryID := c.Param("id")
	var tag barkat.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if httpErr := h.tagMgr.CreateTag(c.Request.Context(), entryID, &tag); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, tag)
}

func (h *TagHandlerImpl) HandleListTags(c *gin.Context) {
	entryID := c.Param("id")
	tagType := c.Query("type")
	tags, httpErr := h.tagMgr.ListTags(c.Request.Context(), entryID, tagType)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

func (h *TagHandlerImpl) HandleDeleteTag(c *gin.Context) {
	entryID := c.Param("id")
	tagID := c.Param("tagId")
	if httpErr := h.tagMgr.DeleteTag(c.Request.Context(), entryID, tagID); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.Status(http.StatusNoContent)
}
