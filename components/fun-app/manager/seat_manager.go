package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

// SeatManagerInterface handles seat-related saga processing and publishing.
// PublishAllocateSeat emits AllocateSeat command (C2) downstream.
// AllocateSeat decides reservation vs waitlist and emits the appropriate seat event;
// on capacity-unavailable it returns an error to trigger middleware retry.
type SeatManagerInterface interface {
	PublishAllocateSeat(ctx context.Context, enrollment fun.Enrollment) common.HttpError
	AllocateSeat(ctx context.Context, cmd fun.AllocateSeatCmdV1) common.HttpError
}

type SeatManager struct {
	SeatPublisher publisher.SeatAllocationPublisher
}

const (
	seatWaitlistThreshold        = 5 // TODO: move to config when real capacity is implemented
	seatWaitlistedReasonCapacity = "capacity_unavailable"
)

// NewSeatManager constructs a seat-only manager that publishes seat events.
func NewSeatManager(seatPublisher publisher.SeatAllocationPublisher) SeatManagerInterface {
	return &SeatManager{
		SeatPublisher: seatPublisher,
	}
}

// PublishAllocateSeat emits the AllocateSeat command for async processing.
func (sm *SeatManager) PublishAllocateSeat(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	return sm.SeatPublisher.AllocateSeat(ctx, enrollment)
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
