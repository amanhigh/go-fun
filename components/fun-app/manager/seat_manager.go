package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

// SeatManagerInterface handles seat-related saga processing and publishing.
// PublishAllocateSeat emits AllocateSeat command (C2) downstream.
// AllocateSeat decides reservation vs waitlist and emits the appropriate seat event.
// Capacity-unavailable produces a successful SeatWaitlisted event; only technical
// errors from publishing or processing are returned as retryable errors.
type SeatManagerInterface interface {
	PublishAllocateSeat(ctx context.Context, enrollment fun.Enrollment) common.HttpError
	AllocateSeat(ctx context.Context, cmd fun.AllocateSeatCmdV1) common.HttpError
	PublishSeatAllocationFailed(ctx context.Context, enrollment fun.Enrollment, reason string) common.HttpError
}

type SeatManager struct {
	SeatPublisher publisher.SeatAllocationPublisher
}

const (
	seatWaitlistThreshold        = 5 // TODO: move to config when real capacity is implemented
	seatWaitlistedReasonCapacity = "capacity_unavailable"
)

// NewSeatManager constructs a seat-only manager that publishes seat events.
func NewSeatManager(seatPublisher publisher.SeatAllocationPublisher) *SeatManager {
	return &SeatManager{
		SeatPublisher: seatPublisher,
	}
}

var _ SeatManagerInterface = (*SeatManager)(nil)

// PublishAllocateSeat emits the AllocateSeat command for async processing.
func (sm *SeatManager) PublishAllocateSeat(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	return sm.SeatPublisher.AllocateSeat(ctx, enrollment)
}

// PublishSeatAllocationFailed emits the terminal seat-allocation failure event.
func (sm *SeatManager) PublishSeatAllocationFailed(ctx context.Context, enrollment fun.Enrollment, reason string) common.HttpError {
	return sm.SeatPublisher.SeatAllocationFailed(ctx, enrollment, reason)
}

// AllocateSeat processes AllocateSeat command and emits SeatReserved or SeatWaitlisted.
// On technical failure it returns an error; waitlist is not a failure.
// No DB writes here; persistence happens in subsequent event handlers.
func (sm *SeatManager) AllocateSeat(ctx context.Context, cmd fun.AllocateSeatCmdV1) common.HttpError {
	enrollment := fun.Enrollment{ID: cmd.EnrollmentID, PersonID: cmd.PersonID, Grade: cmd.Grade}
	if cmd.Grade >= seatWaitlistThreshold {
		return sm.SeatPublisher.SeatWaitlisted(ctx, enrollment, seatWaitlistedReasonCapacity)
	}
	return sm.SeatPublisher.SeatReserved(ctx, enrollment)
}
