package fun

// messaging.go is the single SAGA messaging vocabulary for the enrollment flow.
//
// Every Watermill topic, metadata key, and domain-level constant that
// participates in message routing lives here.

// ---------------------------------------------------------------------------
// Command topics – commands are published to trigger an action.
// ---------------------------------------------------------------------------

const (
	// TopicEnrollCmd kicks off the enrollment saga.
	TopicEnrollCmd = "funapp.enrollment.command.enroll.v1"
	// TopicAllocateSeatCmd requests seat allocation for an enrollment.
	TopicAllocateSeatCmd = "funapp.enrollment.command.allocate_seat.v1"
)

// ---------------------------------------------------------------------------
// Event topics – events announce something that already happened.
// ---------------------------------------------------------------------------

const (
	// TopicSeatReservedEvt signals that a seat was successfully reserved.
	TopicSeatReservedEvt = "funapp.enrollment.event.seat_reserved.v1"
	// TopicSeatWaitlistedEvt signals that the enrollment was waitlisted.
	TopicSeatWaitlistedEvt = "funapp.enrollment.event.seat_waitlisted.v1"
	// TopicSeatAllocationFailedEvt signals that seat allocation exhausted retries.
	TopicSeatAllocationFailedEvt = "funapp.enrollment.event.seat_allocation_failed.v1"
)
