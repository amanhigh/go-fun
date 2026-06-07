package barkat

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

const (
	// AuditBase is the base route for audit framework APIs.
	AuditBase = common.APIV1 + "/audits"
)

// AuditID is a typed identifier for audit plugins.
// The zero value is invalid; use predefined constants.
type AuditID string

const (
	// AuditIDAlertCoverage identifies the Alert Coverage audit plugin.
	AuditIDAlertCoverage AuditID = "alert-coverage"
	// AuditIDStaleReview identifies the Stale Review audit plugin.
	AuditIDStaleReview AuditID = "stale-review"
)

const (
	// AuditFindingNoAlertTicker identifies tickers without an Alert ticker mapping.
	AuditFindingNoAlertTicker = "NO_ALERT_TICKER"
	// AuditFindingNoAlerts identifies mapped tickers without price alerts.
	AuditFindingNoAlerts = "NO_ALERTS"
	// AuditFindingSingleAlert identifies mapped tickers with only one price alert.
	AuditFindingSingleAlert = "SINGLE_ALERT"
	// AuditFindingStaleTicker identifies tickers not opened within the review window.
	AuditFindingStaleTicker = "STALE_TICKER"
)

const (
	// DefaultStaleReviewThresholdDays is the default number of days after which a ticker is considered stale.
	DefaultStaleReviewThresholdDays = 90
)

const (
	// AuditSeverityMedium indicates an operator-relevant audit warning.
	AuditSeverityMedium = "MEDIUM"
	// AuditSeverityHigh indicates a high-priority audit finding.
	AuditSeverityHigh = "HIGH"
)

// Audit describes one active audit check exposed by the catalog API.
type Audit struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

// AuditCatalog is the response body for listing active audit checks.
type AuditCatalog struct {
	Audits []Audit `json:"audits"`
}

// AuditResult is the response body for one executed audit check.
type AuditResult struct {
	AuditID     string                   `json:"audit_id"`
	GeneratedAt time.Time                `json:"generated_at"`
	Counts      map[string]int           `json:"counts"`
	Findings    []AuditFinding           `json:"findings"`
	Metadata    common.PaginatedResponse `json:"metadata"`
}

// AuditFinding describes one operator-facing audit gap.
type AuditFinding struct {
	Code     string            `json:"code"`
	Target   string            `json:"target"`
	Severity string            `json:"severity"`
	Data     map[string]string `json:"data,omitempty"`
}

// AuditPath binds and validates the :audit-id path parameter.
// Audit ID validation is delegated to the plugin registry.
type AuditPath struct {
	AuditID string `uri:"audit-id" json:"-" binding:"required"`
}

// AlertCoverageRow contains the repository projection used by the Alert Coverage audit.
type AlertCoverageRow struct {
	Ticker           string
	AlertTickerCount int64
	PriceAlertCount  int64
}
