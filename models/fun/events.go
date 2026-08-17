package fun

import "time"

// FIXME: Introduce a base EnrollEvent for fields shared by enrollment events.
type SeatReservedEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId" validate:"required"`
	PersonID     string    `json:"personId" validate:"required"`
	Grade        int       `json:"grade" validate:"required,min=1,max=12"`
	ReservedAt   time.Time `json:"reservedAt" validate:"required"`
}

// HACK: Reorganise around Domain names not Event.
type SeatWaitlistedEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId" validate:"required"`
	PersonID     string    `json:"personId" validate:"required"`
	Grade        int       `json:"grade" validate:"required,min=1,max=12"`
	Reason       string    `json:"reason" validate:"required"`
	WaitlistedAt time.Time `json:"waitlistedAt" validate:"required"`
}

type SeatAllocationFailedEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId" validate:"required"`
	PersonID     string    `json:"personId" validate:"required"`
	Reason       string    `json:"reason" validate:"required"`
	FailedAt     time.Time `json:"failedAt" validate:"required"`
}
