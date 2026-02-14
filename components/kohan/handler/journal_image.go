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
	// HandleCreateImage handles POST /v1/journal-entries/:id/images
	HandleCreateImage(c *gin.Context)
	// HandleListImages handles GET /v1/journal-entries/:id/images
	HandleListImages(c *gin.Context)
	// HandleDeleteImage handles DELETE /v1/journal-entries/:id/images/:imageId
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
	entryID := c.Param("id")
	var image barkat.Image
	if err := c.ShouldBindJSON(&image); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if httpErr := h.imageMgr.CreateImage(c.Request.Context(), entryID, &image); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, image)
}

func (h *ImageHandlerImpl) HandleListImages(c *gin.Context) {
	entryID := c.Param("id")
	images, httpErr := h.imageMgr.ListImages(c.Request.Context(), entryID)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"images": images})
}

func (h *ImageHandlerImpl) HandleDeleteImage(c *gin.Context) {
	entryID := c.Param("id")
	imageID := c.Param("imageId")
	if httpErr := h.imageMgr.DeleteImage(c.Request.Context(), entryID, imageID); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.Status(http.StatusNoContent)
}
