package tax

import (
	"fmt"
	"net/http"
	"time"
)

type closestDateError struct {
	requestedDate time.Time
	closestDate   time.Time
}

func NewClosestDateError(requested, closest time.Time) ClosestDateError {
	return &closestDateError{
		requestedDate: requested,
		closestDate:   closest,
	}
}

func (e *closestDateError) Error() string {
	return fmt.Sprintf("exact rate not found for %v, using closest available date %v",
		e.requestedDate.Format(time.DateOnly),
		e.closestDate.Format(time.DateOnly))
}

func (e *closestDateError) Code() int {
	// Using 200 as this is an expected case
	return http.StatusOK
}

func (e *closestDateError) GetClosestDate() time.Time {
	return e.closestDate
}

func (e *closestDateError) GetRequestedDate() time.Time {
	return e.requestedDate
}
