package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

// SeatManagerInterface handles only seat-related saga processing and publishing.
// - AllocateSeat: decides reservation vs waitlist and emits the appropriate event
// - OnSeatWaitlistedEvt: persists WAITLISTED as a sink (idempotent), no further publishing

// Seat reservation confirmation (after SeatReserved event) is handled by EnrollmentManager via handler.
type SeatManagerInterface interface {
	AllocateSeat(ctx context.Context, cmd fun.AllocateSeatCmdV1) common.HttpError
	OnSeatWaitlistedEvt(ctx context.Context, evt fun.SeatWaitlistedEvtV1) common.HttpError
}

type SeatManager struct {
	SeatPublisher     publisher.SeatAllocationPublisher
	EnrollmentManager EnrollmentManagerInterface
}

const seatWaitlistThreshold = 5 // TODO: move to config when real capacity is implemented

// NewSeatManager constructs a seat-only manager that publishes seat events and
// delegates enrollment status persistence to the EnrollmentManager where needed.
func NewSeatManager(seatPublisher publisher.SeatAllocationPublisher, enrollmentManager EnrollmentManagerInterface) SeatManagerInterface {
	return &SeatManager{
		SeatPublisher:     seatPublisher,
		EnrollmentManager: enrollmentManager,
	}
}

// OnAllocateSeatCmd processes AllocateSeat command and emits SeatReserved or SeatWaitlisted.
// No DB writes here; persistence happens in subsequent event handlers.
func (sm *SeatManager) AllocateSeat(ctx context.Context, cmd fun.AllocateSeatCmdV1) common.HttpError {
	e := fun.Enrollment{ID: cmd.EnrollmentID, PersonID: cmd.PersonID, Grade: cmd.Grade}
	if cmd.Grade >= seatWaitlistThreshold {
		return sm.SeatPublisher.SeatWaitlisted(ctx, e, seatWaitlistedReasonCapacity)
	}
	return sm.SeatPublisher.SeatReserved(ctx, e)
}

// OnSeatWaitlistedEvt idempotently persists WAITLISTED status as a sink.
// No event is (re)published here to avoid loops; EnrollmentManager handles persistence.
func (sm *SeatManager) OnSeatWaitlistedEvt(ctx context.Context, evt fun.SeatWaitlistedEvtV1) common.HttpError {
	e := fun.Enrollment{ID: evt.EnrollmentID, PersonID: evt.PersonID, Grade: evt.Grade}
	return sm.EnrollmentManager.UpdateToWaitlisted(ctx, e)
}
