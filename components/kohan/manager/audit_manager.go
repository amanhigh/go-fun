package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/manager/audit"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// AuditManager provides audit catalog and execution business logic.
type AuditManager interface {
	// ListAudits returns active audit checks in display order.
	ListAudits(ctx context.Context) (barkat.AuditCatalog, common.HttpError)
	// ExecuteAudit runs one audit check and returns paginated findings.
	ExecuteAudit(ctx context.Context, auditID string, query common.Pagination) (barkat.AuditResult, common.HttpError)
}

type AuditManagerImpl struct {
	registry *audit.PluginRegistry
}

var _ AuditManager = (*AuditManagerImpl)(nil)

// NewAuditManager creates a new AuditManager backed by a plugin registry.
func NewAuditManager(registry *audit.PluginRegistry) *AuditManagerImpl {
	return &AuditManagerImpl{registry: registry}
}

func (m *AuditManagerImpl) ListAudits(_ context.Context) (barkat.AuditCatalog, common.HttpError) {
	return barkat.AuditCatalog{Audits: m.registry.ListCatalog()}, nil
}

func (m *AuditManagerImpl) ExecuteAudit(ctx context.Context, auditID string, query common.Pagination) (barkat.AuditResult, common.HttpError) {
	plugin, httpErr := m.registry.GetPlugin(barkat.AuditID(auditID))
	if httpErr != nil {
		return barkat.AuditResult{}, httpErr
	}
	return plugin.Execute(ctx, query)
}
