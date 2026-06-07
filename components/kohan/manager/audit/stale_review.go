package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// StaleReviewPlugin detects tracked instruments that have not been opened
// within the configured review window.
type StaleReviewPlugin struct {
	repo          repository.AuditRepository
	thresholdDays int
}

// Compile-time interface check.
var _ Plugin = (*StaleReviewPlugin)(nil)

// NewStaleReviewPlugin creates a new StaleReviewPlugin with the given threshold.
func NewStaleReviewPlugin(repo repository.AuditRepository) *StaleReviewPlugin {
	return &StaleReviewPlugin{
		repo:          repo,
		thresholdDays: barkat.DefaultStaleReviewThresholdDays,
	}
}

func (p *StaleReviewPlugin) ID() barkat.AuditID {
	return barkat.AuditIDStaleReview
}

func (p *StaleReviewPlugin) Title() string {
	return "Stale Review"
}

func (p *StaleReviewPlugin) Description() string {
	return fmt.Sprintf("Tracked instruments that have not been opened within the last %d days.", p.thresholdDays)
}

func (p *StaleReviewPlugin) Order() int {
	return 2
}

func (p *StaleReviewPlugin) Execute(ctx context.Context, query common.Pagination) (barkat.AuditResult, common.HttpError) {
	var result barkat.AuditResult
	err := p.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		cutoff := time.Now().UTC().AddDate(0, 0, -p.thresholdDays)
		rows, httpErr := p.repo.ListStaleReviewTickers(c, cutoff)
		if httpErr != nil {
			return httpErr
		}

		findings := buildStaleReviewFindings(rows)
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

func buildStaleReviewFindings(rows []barkat.Ticker) []barkat.AuditFinding {
	findings := make([]barkat.AuditFinding, 0, len(rows))
	for _, row := range rows {
		findings = append(findings, buildStaleReviewFinding(row))
	}
	return findings
}

func buildStaleReviewFinding(ticker barkat.Ticker) barkat.AuditFinding {
	return barkat.AuditFinding{
		Code:     barkat.AuditFindingStaleTicker,
		Target:   ticker.Ticker,
		Severity: barkat.AuditSeverityMedium,
		Data: map[string]string{
			"last_opened_at": ticker.LastOpenedAt.Format(time.RFC3339),
		},
	}
}
