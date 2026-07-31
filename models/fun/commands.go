package fun

import "time"

// EnrollCmdV1 triggers the enrollment saga flow.
type EnrollCmdV1 struct {
	EnrollmentID string    `json:"enrollmentId" validate:"required"`
	PersonID     string    `json:"personId" validate:"required"`
	Grade        int       `json:"grade" validate:"required,min=1,max=12"`
	Status       string    `json:"status" validate:"required,eq=SEAT_ALLOCATION_INITIATED"`
	RequestedAt  time.Time `json:"requestedAt" validate:"required"`
}

// AllocateSeatCmdV1 requests seat allocation for an enrollment.
type AllocateSeatCmdV1 struct {
	EnrollmentID string    `json:"enrollmentId" validate:"required"`
	PersonID     string    `json:"personId" validate:"required"`
	Grade        int       `json:"grade" validate:"required,min=1,max=12"`
	RequestedAt  time.Time `json:"requestedAt" validate:"required"`
}
