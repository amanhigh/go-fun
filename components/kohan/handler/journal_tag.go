package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/gin-gonic/gin"
)

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
