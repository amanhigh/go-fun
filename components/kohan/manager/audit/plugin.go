// Package audit defines the audit plugin contract used by the audit framework.
// Each plugin represents one independently verifiable audit check.
package audit

import (
	"context"

	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// Plugin is the interface that every audit plugin must implement.
// Plugins provide their own identity, catalog metadata, and execution logic.
// This mirrors the frontend BasePlugin pattern.
type Plugin interface {
	// ID returns the unique plugin identifier matching the `oneof` validation tag.
	ID() barkat.AuditID

	// Title returns the human-readable display name.
	Title() string

	// Description returns the operator-facing description or tooltip text.
	Description() string

	// Order returns the display order for the audit panel (lower = earlier).
	Order() int

	// Execute runs the audit check and returns paginated findings.
	// The implementation must set AuditID in the result to its own ID.
	Execute(ctx context.Context, query common.Pagination) (barkat.AuditResult, common.HttpError)
}

// Metadata returns the catalog Audit entry for this plugin.
func Metadata(p Plugin) barkat.Audit {
	return barkat.Audit{
		ID:          string(p.ID()),
		Title:       p.Title(),
		Description: p.Description(),
		Order:       p.Order(),
	}
}
