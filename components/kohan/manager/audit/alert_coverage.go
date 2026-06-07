package audit

import (
	"context"
	"strconv"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// AlertCoveragePlugin checks tracked instruments for missing or insufficient
// price-alert coverage. It evaluates active READY-state tickers and skips
// actively watched (WATCHED) and blacklisted (BLACKLIST) instruments.
type AlertCoveragePlugin struct {
	repo repository.AuditRepository
}

// Compile-time interface check.
var _ Plugin = (*AlertCoveragePlugin)(nil)

// NewAlertCoveragePlugin creates a new AlertCoveragePlugin.
func NewAlertCoveragePlugin(repo repository.AuditRepository) *AlertCoveragePlugin {
	return &AlertCoveragePlugin{repo: repo}
}

func (p *AlertCoveragePlugin) ID() barkat.AuditID {
	return barkat.AuditIDAlertCoverage
}

func (p *AlertCoveragePlugin) Title() string {
	return "Alert Coverage"
}

func (p *AlertCoveragePlugin) Description() string {
	return "Tracked instruments with missing or insufficient price-alert coverage."
}

func (p *AlertCoveragePlugin) Order() int {
	return 1
}

func (p *AlertCoveragePlugin) Execute(ctx context.Context, query common.Pagination) (barkat.AuditResult, common.HttpError) {
	var result barkat.AuditResult
	err := p.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		rows, httpErr := p.repo.ListAlertCoverageRows(c)
		if httpErr != nil {
			return httpErr
		}

		findings := buildAlertCoverageFindings(rows)
		result = barkat.AuditResult{
			AuditID:     string(p.ID()),
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
