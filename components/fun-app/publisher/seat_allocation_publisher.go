package publisher

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

// SeatAllocationPublisher handles seat-related saga events.
type SeatAllocationPublisher interface {
	AllocateSeat(ctx context.Context, enrollment fun.Enrollment) common.HttpError
	SeatReserved(ctx context.Context, enrollment fun.Enrollment) common.HttpError
	SeatWaitlisted(ctx context.Context, enrollment fun.Enrollment, reason string) common.HttpError
	SeatAllocationFailed(ctx context.Context, enrollment fun.Enrollment, reason string) common.HttpError
}

type seatAllocationPublisher struct {
	base BasePublisher
}

// NewSeatAllocationPublisher constructs a SeatAllocationPublisher.
func NewSeatAllocationPublisher(base BasePublisher) SeatAllocationPublisher {
	return &seatAllocationPublisher{base: base}
}

func (sap *seatAllocationPublisher) AllocateSeat(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	payload := fun.AllocateSeatCmdV1{
		EnrollmentID: enrollment.ID,
		PersonID:     enrollment.PersonID,
		Grade:        enrollment.Grade,
		RequestedAt:  time.Now().UTC(),
	}

	return sap.base.PublishChild(ctx, fun.TopicAllocateSeatCmd, payload)
}

func (sap *seatAllocationPublisher) SeatReserved(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	payload := fun.SeatReservedEvtV1{
		EnrollmentID: enrollment.ID,
		PersonID:     enrollment.PersonID,
		Grade:        enrollment.Grade,
		ReservedAt:   time.Now().UTC(),
	}

	return sap.base.PublishChild(ctx, fun.TopicSeatReservedEvt, payload)
}

func (sap *seatAllocationPublisher) SeatWaitlisted(ctx context.Context, enrollment fun.Enrollment, reason string) common.HttpError {
	payload := fun.SeatWaitlistedEvtV1{
		EnrollmentID: enrollment.ID,
		PersonID:     enrollment.PersonID,
		Grade:        enrollment.Grade,
		Reason:       reason,
		WaitlistedAt: time.Now().UTC(),
	}

	return sap.base.PublishChild(ctx, fun.TopicSeatWaitlistedEvt, payload)
}

func (sap *seatAllocationPublisher) SeatAllocationFailed(ctx context.Context, enrollment fun.Enrollment, reason string) common.HttpError {
	payload := fun.SeatAllocationFailedEvtV1{
		EnrollmentID: enrollment.ID,
		PersonID:     enrollment.PersonID,
		Reason:       reason,
		FailedAt:     time.Now().UTC(),
	}

	return sap.base.PublishChild(ctx, fun.TopicSeatAllocationFailedEvt, payload)
}
