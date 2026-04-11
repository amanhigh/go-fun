package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/ui/pages"
	"github.com/gin-gonic/gin"
)

// JournalPortal handles journal UI portal routes.
type JournalPortal interface {
	// ListJournals renders the journal list page.
	ListJournals(ctx *gin.Context)
	// DisplayJournal renders the dedicated journal detail page.
	DisplayJournal(ctx *gin.Context)
	// ImagePath returns the base path for journal images.
	ImagePath() string
}

type JournalDetailPath struct {
	JournalID string `uri:"id" binding:"required,journal_id"`
}

type JournalPortalImpl struct {
	imagePath string
}

// NewJournalPortal creates a new JournalPortal.
func NewJournalPortal(imagePath string) *JournalPortalImpl {
	return &JournalPortalImpl{imagePath: imagePath}
}

var _ JournalPortal = (*JournalPortalImpl)(nil)

func (h *JournalPortalImpl) ImagePath() string {
	return h.imagePath
}

func (h *JournalPortalImpl) ListJournals(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")
	if err := pages.JournalPage().Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to render journal page"})
	}
}

func (h *JournalPortalImpl) DisplayJournal(ctx *gin.Context) {
	var path JournalDetailPath
	if bindErr := ctx.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		ctx.JSON(httpErr.Code(), httpErr)
		return
	}

	ctx.Header("Content-Type", "text/html")
	if err := pages.JournalDetailPage(path.JournalID).Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to render journal detail page"})
	}
}
