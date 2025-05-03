package tax

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

// Base interface for all CSV records
type CSVRecord interface {
	GetKey() string                               // For key-based lookups (symbol/date)
	GetDate() (time.Time, common.HttpError)       // For date-based operations, returns error if parsing fails
	IsValid() bool                                // Validate record fields
}

// Exchangeable now extends CSVRecord
type Exchangeable interface {
	CSVRecord
	GetUSDAmount() float64
	SetTTRate(float64)
	SetTTDate(time.Time)
}
