package fun

import "time"

// AllocateSeatCmdV1 requests seat allocation for an enrollment.
type AllocateSeatCmdV1 struct {
	EnrollmentID string    `json:"enrollmentId" validate:"required"`
	PersonID     string    `json:"personId" validate:"required"`
	Grade        int       `json:"grade" validate:"required,min=1,max=12"`
	RequestedAt  time.Time `json:"requestedAt" validate:"required"`
}

// SeatReservedEvtV1 signals that a seat was successfully reserved.
type SeatReservedEvtV1 struct {
	EnrollmentEvent
	Grade      int       `json:"grade" validate:"required,min=1,max=12"`
	ReservedAt time.Time `json:"reservedAt" validate:"required"`
}

// SeatWaitlistedEvtV1 signals that the enrollment was waitlisted.
type SeatWaitlistedEvtV1 struct {
	EnrollmentEvent
	Grade        int       `json:"grade" validate:"required,min=1,max=12"`
	Reason       string    `json:"reason" validate:"required"`
	WaitlistedAt time.Time `json:"waitlistedAt" validate:"required"`
}

// SeatAllocationFailedEvtV1 signals that seat allocation exhausted retries.
type SeatAllocationFailedEvtV1 struct {
	EnrollmentEvent
	Reason   string    `json:"reason" validate:"required"`
	FailedAt time.Time `json:"failedAt" validate:"required"`
}
