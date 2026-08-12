package fun

import "time"

type SeatReservedEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId"`
	StudentID     string    `json:"studentId"`
	Grade        int       `json:"grade"`
	ReservedAt   time.Time `json:"reservedAt"`
}

type SeatWaitlistedEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId"`
	StudentID     string    `json:"studentId"`
	Grade        int       `json:"grade"`
	Reason       string    `json:"reason"`
	WaitlistedAt time.Time `json:"waitlistedAt"`
}

type EnrollmentConfirmedEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId"`
	StudentID     string    `json:"studentId"`
	ConfirmedAt  time.Time `json:"confirmedAt"`
}

type EnrollmentCancelledEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId"`
	StudentID     string    `json:"studentId"`
	Reason       string    `json:"reason"`
	CancelledAt  time.Time `json:"cancelledAt"`
}
