package common

// Common type for context keys used across packages
type ContextKey string

// Common date formats used across application
const (
	DateOnly = "2006-01-02"          // For dates without time
	DateTime = "2006-01-02 15:04:05" // For dates with time
)
