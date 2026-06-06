package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
)

// AuditHandler provides HTTP handlers for audit framework operations.
type AuditHandler interface {
	HandleListAudits(c *gin.Context)
	HandleExecuteAudit(c *gin.Context)
}

type AuditHandlerImpl struct {
	auditManager manager.AuditManager
}

var _ AuditHandler = (*AuditHandlerImpl)(nil)

// NewAuditHandler creates a new AuditHandlerImpl.
func NewAuditHandler(auditManager manager.AuditManager) *AuditHandlerImpl {
	return &AuditHandlerImpl{auditManager: auditManager}
}

func (h *AuditHandlerImpl) HandleListAudits(c *gin.Context) {
	response, httpErr := h.auditManager.ListAudits(c.Request.Context())
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(response))
}

func (h *AuditHandlerImpl) HandleExecuteAudit(c *gin.Context) {
	auditID := c.Param("audit-id")
	if auditID == "" {
		c.JSON(http.StatusNotFound, common.NewFailEnvelope(map[string]string{"audit-id": "Audit not found"}))
		return
	}

	var query common.Pagination
	if bindErr := c.ShouldBindQuery(&query); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	response, httpErr := h.auditManager.ExecuteAudit(c.Request.Context(), auditID, query)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(response))
}
