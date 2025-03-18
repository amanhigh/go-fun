package tax

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

// Base interface for all CSV records
type CSVRecord interface {
	GetKey() string // For key-based lookups (symbol/date)
	// FIXME: Should we Return error to caller?
	GetDate() time.Time // For date-based operations
	IsValid() bool      // Validate record fields
}

// Exchangeable now extends CSVRecord
type Exchangeable interface {
	CSVRecord
	GetUSDAmount() float64
	SetTTRate(float64)
	SetTTDate(time.Time)
}

type ClosestDateError interface {
	common.HttpError
	GetClosestDate() time.Time
	GetRequestedDate() time.Time
}
