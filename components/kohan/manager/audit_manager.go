package manager

import (
	"context"
	"strconv"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/repository"
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
	repo repository.AuditRepository
}

var _ AuditManager = (*AuditManagerImpl)(nil)

// NewAuditManager creates a new AuditManager.
func NewAuditManager(repo repository.AuditRepository) *AuditManagerImpl {
	return &AuditManagerImpl{repo: repo}
}

func (m *AuditManagerImpl) ListAudits(_ context.Context) (barkat.AuditCatalog, common.HttpError) {
	return barkat.AuditCatalog{Audits: []barkat.Audit{alertCoverageAudit()}}, nil
}

func (m *AuditManagerImpl) ExecuteAudit(ctx context.Context, auditID string, query common.Pagination) (barkat.AuditResult, common.HttpError) {
	return m.executeAlertCoverageAudit(ctx, query)
}

func (m *AuditManagerImpl) executeAlertCoverageAudit(ctx context.Context, query common.Pagination) (barkat.AuditResult, common.HttpError) {
	var result barkat.AuditResult
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		rows, httpErr := m.repo.ListAlertCoverageRows(c)
		if httpErr != nil {
			return httpErr
		}

		findings := buildAlertCoverageFindings(rows)
		result = barkat.AuditResult{
			AuditID:     string(barkat.AuditIDAlertCoverage),
			GeneratedAt: time.Now().UTC(),
			Counts:      countAuditFindings(findings),
			Findings:    paginateAuditFindings(findings, query),
			Metadata: common.PaginatedResponse{
				Total:  int64(len(findings)),
				Offset: query.Offset,
				Limit:  query.Limit,
			},
		}
		return nil
	})
	return result, err
}

func alertCoverageAudit() barkat.Audit {
	return barkat.Audit{
		ID:          string(barkat.AuditIDAlertCoverage),
		Title:       "Alert Coverage",
		Description: "Tracked instruments with missing or insufficient price-alert coverage.",
		Order:       1,
	}
}

func buildAlertCoverageFindings(rows []barkat.AlertCoverageRow) []barkat.AuditFinding {
	findings := make([]barkat.AuditFinding, 0, len(rows))
	for _, row := range rows {
		finding, ok := buildAlertCoverageFinding(row)
		if ok {
			findings = append(findings, finding)
		}
	}
	return findings
}

func buildAlertCoverageFinding(row barkat.AlertCoverageRow) (barkat.AuditFinding, bool) {
	data := map[string]string{
		"alert_ticker_count": strconv.FormatInt(row.AlertTickerCount, 10),
		"price_alert_count":  strconv.FormatInt(row.PriceAlertCount, 10),
	}

	if row.AlertTickerCount == 0 {
		return barkat.AuditFinding{Code: barkat.AuditFindingNoAlertTicker, Target: row.Ticker, Severity: barkat.AuditSeverityHigh, Data: data}, true
	}
	if row.PriceAlertCount == 0 {
		return barkat.AuditFinding{Code: barkat.AuditFindingNoAlerts, Target: row.Ticker, Severity: barkat.AuditSeverityMedium, Data: data}, true
	}
	if row.PriceAlertCount == 1 {
		return barkat.AuditFinding{Code: barkat.AuditFindingSingleAlert, Target: row.Ticker, Severity: barkat.AuditSeverityHigh, Data: data}, true
	}
	return barkat.AuditFinding{}, false
}

func countAuditFindings(findings []barkat.AuditFinding) map[string]int {
	counts := make(map[string]int)
	for _, finding := range findings {
		counts[finding.Code]++
	}
	return counts
}

func paginateAuditFindings(findings []barkat.AuditFinding, query common.Pagination) []barkat.AuditFinding {
	if query.Offset >= len(findings) {
		return []barkat.AuditFinding{}
	}
	end := min(query.Offset+query.Limit, len(findings))
	return findings[query.Offset:end]
}
