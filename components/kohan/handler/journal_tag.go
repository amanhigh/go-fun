//nolint:dupl // Intentional CRUD pattern shared with NoteHandler for different sub-resources
package handler

// TagHandler provides HTTP handlers for journal tag operations.
// Tags categorize entries with reason codes (oe, dep, etc.) or management labels.

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
)

// TagHandler provides HTTP handlers for journal tag operations.
type TagHandler interface {
	HandleCreateTag(c *gin.Context)
	HandleListTags(c *gin.Context)
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
	var path barkat.JournalPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	var tag barkat.Tag
	if bindErr := c.ShouldBindJSON(&tag); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	createdTag, httpErr := h.tagMgr.CreateTag(c.Request.Context(), path.JournalID, tag)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, common.NewEnvelope(createdTag))
}

func (h *TagHandlerImpl) HandleListTags(c *gin.Context) {
	var path barkat.JournalPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	query := barkat.TagQuery{}

	if bindErr := c.ShouldBindQuery(&query); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	tagList, httpErr := h.tagMgr.ListTags(c.Request.Context(), path.JournalID, query.Type)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(tagList))
}

func (h *TagHandlerImpl) HandleDeleteTag(c *gin.Context) {
	var path barkat.TagPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	if httpErr := h.tagMgr.DeleteTag(c.Request.Context(), path.JournalID, path.TagID); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.Status(http.StatusNoContent)
}
