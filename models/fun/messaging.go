package fun

// messaging.go is the single SAGA messaging vocabulary for the enrollment flow.
//
// Every Watermill topic, metadata key, router handler identifier, and
// domain-level constant that participates in message routing lives here.
// Router handler IDs are kept in this catalog for learning visibility.

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
	// TopicEnrollmentConfirmedEvt signals that the enrollment is confirmed.
	TopicEnrollmentConfirmedEvt = "funapp.enrollment.event.enrollment_confirmed.v1"
	// TopicEnrollmentCancelledEvt signals that the enrollment was cancelled.
	TopicEnrollmentCancelledEvt = "funapp.enrollment.event.enrollment_cancelled.v1"
)

// ---------------------------------------------------------------------------
// Poison / dead-letter topics – messages land here after retry exhaustion.
// ---------------------------------------------------------------------------

const (
	// TopicPoisonedAllocateSeatCmd is the allocation-specific poison topic
	// used when an AllocateSeatCmd fails after all retries.
	TopicPoisonedAllocateSeatCmd = "funapp.enrollment.poison.allocate_seat.v1"
)

// ---------------------------------------------------------------------------
// Watermill router handler identifiers.
//
// These constants are used in router handler registration (AddConsumerHandler /
// AddNoPublisherHandler) and provide a central catalog for visibility.
// ---------------------------------------------------------------------------

const (
	RouterHandlerIDRequestSeatAllocation           = "request_seat_allocation"
	RouterHandlerIDAllocateSeat                    = "allocate_seat"
	RouterHandlerIDConfirmEnrollment               = "confirm_enrollment"
	RouterHandlerIDWaitlistEnrollment              = "waitlist_enrollment"
	RouterHandlerIDRecordEnrollmentConfirmation    = "record_enrollment_confirmation"
	RouterHandlerIDRecordEnrollmentCancellation    = "record_enrollment_cancellation"
	RouterHandlerIDCompensateSeatAllocationFailure = "compensate_seat_allocation_failure"
)

// ---------------------------------------------------------------------------
// Message metadata keys – attached as Watermill message metadata headers.
// ---------------------------------------------------------------------------

const (
	MetadataEnrollmentID = "enrollment_id"
	MetadataPersonID     = "person_id"
)

// ---------------------------------------------------------------------------
// Domain constants
// ---------------------------------------------------------------------------

const (
	// EnrollmentCancellationReasonSeatAllocationFailed is the reason stored in
	// EnrollmentCancelledEvtV1 when the cancellation is triggered by an
	// allocation failure after retries are exhausted.
	EnrollmentCancellationReasonSeatAllocationFailed = "seat_allocation_failed_after_retries"
)
