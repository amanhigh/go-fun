package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/gin-gonic/gin"
)

// ImageHandler provides HTTP handlers for journal image operations.
type ImageHandler interface {
	HandleCreateImage(c *gin.Context)
	HandleListImages(c *gin.Context)
	HandleDeleteImage(c *gin.Context)
}

type ImageHandlerImpl struct {
	imageMgr manager.ImageManager
}

var _ ImageHandler = (*ImageHandlerImpl)(nil)

// NewImageHandler creates a new ImageHandler.
func NewImageHandler(imageMgr manager.ImageManager) *ImageHandlerImpl {
	return &ImageHandlerImpl{imageMgr: imageMgr}
}

func (h *ImageHandlerImpl) HandleCreateImage(c *gin.Context) {
	journalID := c.Param("id")
	var image barkat.Image
	if bindErr := c.ShouldBindJSON(&image); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	createdImage, httpErr := h.imageMgr.CreateImage(c.Request.Context(), journalID, image)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, createdImage)
}

func (h *ImageHandlerImpl) HandleListImages(c *gin.Context) {
	journalID := c.Param("id")
	images, httpErr := h.imageMgr.ListImages(c.Request.Context(), journalID)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"images": images})
}

func (h *ImageHandlerImpl) HandleDeleteImage(c *gin.Context) {
	journalID := c.Param("id")
	imageID := c.Param("imageId")
	if httpErr := h.imageMgr.DeleteImage(c.Request.Context(), journalID, imageID); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.Status(http.StatusNoContent)
}
