package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IndexHandler handles index page HTTP routes.
type IndexHandler interface {
	// HandleIndex renders the index page.
	HandleIndex(c *gin.Context)
}

type IndexHandlerImpl struct{}

func NewIndexHandler() *IndexHandlerImpl {
	return &IndexHandlerImpl{}
}

var _ IndexHandler = (*IndexHandlerImpl)(nil)

func (h *IndexHandlerImpl) HandleIndex(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	if err := IndexPage().Render(c.Request.Context(), c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render"})
	}
}
