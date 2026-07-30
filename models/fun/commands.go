package fun

import "time"

// EnrollCmdV1 triggers the enrollment saga flow.
type EnrollCmdV1 struct {
	EnrollmentID string    `json:"enrollmentId"`
	PersonID     string    `json:"personId"`
	Grade        int       `json:"grade"`
	Status       string    `json:"status"`
	RequestedAt  time.Time `json:"requestedAt"`
}

// AllocateSeatCmdV1 requests seat allocation for an enrollment.
type AllocateSeatCmdV1 struct {
	EnrollmentID string    `json:"enrollmentId"`
	PersonID     string    `json:"personId"`
	Grade        int       `json:"grade"`
	RequestedAt  time.Time `json:"requestedAt"`
}
