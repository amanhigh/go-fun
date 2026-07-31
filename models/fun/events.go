package fun

import "time"

// FIXME: Introduce a base EnrollEvent for fields shared by enrollment events.
type SeatReservedEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId" validate:"required"`
	PersonID     string    `json:"personId" validate:"required"`
	Grade        int       `json:"grade" validate:"required,min=1,max=12"`
	ReservedAt   time.Time `json:"reservedAt" validate:"required"`
}

// TODO: Reorganise around Domain names not Event.
type SeatWaitlistedEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId"`
	PersonID     string    `json:"personId"`
	Grade        int       `json:"grade"`
	Reason       string    `json:"reason"`
	WaitlistedAt time.Time `json:"waitlistedAt"`
}

type EnrollmentConfirmedEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId"`
	PersonID     string    `json:"personId"`
	ConfirmedAt  time.Time `json:"confirmedAt"`
}

type EnrollmentCancelledEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId"`
	PersonID     string    `json:"personId"`
	Reason       string    `json:"reason"`
	CancelledAt  time.Time `json:"cancelledAt"`
}

type SeatAllocationFailedEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId"`
	PersonID     string    `json:"personId"`
	Reason       string    `json:"reason"`
	FailedAt     time.Time `json:"failedAt"`
}
