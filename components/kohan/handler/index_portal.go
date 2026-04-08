package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/components/kohan/ui/pages"
	"github.com/gin-gonic/gin"
)

// IndexPortal handles UI portal routes.
type IndexPortal interface {
	// HandleIndex renders the index portal page.
	HandleIndex(ctx *gin.Context)
}

type IndexPortalImpl struct{}

// NewIndexPortal creates a new IndexPortal.
func NewIndexPortal() *IndexPortalImpl {
	return &IndexPortalImpl{}
}

var _ IndexPortal = (*IndexPortalImpl)(nil)

func (h *IndexPortalImpl) HandleIndex(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")
	if err := pages.IndexPage().Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to render index page"})
	}
}
