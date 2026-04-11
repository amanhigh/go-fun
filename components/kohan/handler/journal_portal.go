package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/components/kohan/ui/pages"
	"github.com/gin-gonic/gin"
)

// JournalPortal handles journal UI portal routes.
type JournalPortal interface {
	// HandleJournal renders the journal portal page.
	HandleJournal(ctx *gin.Context)
}

type JournalPortalImpl struct{}

// NewJournalPortal creates a new JournalPortal.
func NewJournalPortal() *JournalPortalImpl {
	return &JournalPortalImpl{}
}

var _ JournalPortal = (*JournalPortalImpl)(nil)

func (h *JournalPortalImpl) HandleJournal(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")
	if err := pages.JournalPage().Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to render journal page"})
	}
}
