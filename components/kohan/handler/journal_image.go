package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/gin-gonic/gin"
)

// HACK: Image,Notes,Tags should have their own Interface and Impl not in JournalHandler for Handler,Manager,Repo
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
