package tax

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

// SecurityInfo represents basic security information from Yahoo Finance.
// This is the minimal model for Phase 1 covering symbol, name, exchange, and type.
type SecurityInfo struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Type     string `json:"type"`
}

// SplitInfo represents a stock split event.
type SplitInfo struct {
	Date        int64   `json:"date"`
	Numerator   float64 `json:"numerator"`
	Denominator float64 `json:"denominator"`
}

// Validate checks that a SplitInfo has valid event data:
// positive finite numerator and denominator, and a valid usable Unix timestamp.
// The error message includes the ticker and event timestamp for traceability.
func (s SplitInfo) Validate(ticker string) common.HttpError {
	prefix := fmt.Sprintf("ticker %s: split event timestamp %d", ticker, s.Date)
	if s.Date <= 0 {
		return common.NewHttpError(fmt.Sprintf("%s: non-positive timestamp %d", prefix, s.Date), http.StatusBadRequest)
	}
	if s.Numerator <= 0 || math.IsInf(s.Numerator, 0) || math.IsNaN(s.Numerator) {
		return common.NewHttpError(fmt.Sprintf("%s: non-positive or non-finite numerator %f", prefix, s.Numerator), http.StatusBadRequest)
	}
	if s.Denominator <= 0 || math.IsInf(s.Denominator, 0) || math.IsNaN(s.Denominator) {
		return common.NewHttpError(fmt.Sprintf("%s: non-positive or non-finite denominator %f", prefix, s.Denominator), http.StatusBadRequest)
	}
	return nil
}

// Ratio returns the split ratio (Numerator / Denominator).
// Callers must call Validate first — the result is undefined for
// unvalidated or invalid splits with zero/infinity/NaN components.
func (s SplitInfo) Ratio() float64 {
	return s.Numerator / s.Denominator
}

// EffectiveDate returns the UTC calendar date of the split at midnight,
// normalizing the Yahoo Unix timestamp to the start of its UTC day.
func (s SplitInfo) EffectiveDate() time.Time {
	return time.Unix(s.Date, 0).UTC().Truncate(24 * time.Hour) //nolint:mnd
}
