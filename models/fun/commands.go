package fun

import "time"

const (
	TopicEnrollCmd       = "funapp.enrollment.command.enroll.v1"
	TopicAllocateSeatCmd = "funapp.enrollment.command.allocate_seat.v1"

	MetadataEnrollmentID = "enrollment_id"
	MetadataPersonID     = "person_id"
	MetadataCorrelation  = "correlation_id"
	MetadataCausation    = "causation_id"
)

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
