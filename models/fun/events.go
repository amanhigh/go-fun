package fun

import "time"

const (
	TopicSeatReservedEvt           = "funapp.enrollment.event.seat_reserved.v1"
	TopicSeatWaitlistedEvt         = "funapp.enrollment.event.seat_waitlisted.v1"
	TopicEnrollmentConfirmedEvt    = "funapp.enrollment.event.enrollment_confirmed.v1"
	TopicEnrollmentCancelledEvt    = "funapp.enrollment.event.enrollment_cancelled.v1"
	TopicEnrollmentStateTransition = "funapp.enrollment.event.state_transition.v1"
	// Poison queue for failed messages after retries are exhausted
	TopicPoison = "funapp.enrollment.poison"
)

type SeatReservedEvtV1 struct {
	EnrollmentID string    `json:"enrollmentId"`
	PersonID     string    `json:"personId"`
	Grade        int       `json:"grade"`
	ReservedAt   time.Time `json:"reservedAt"`
}

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
