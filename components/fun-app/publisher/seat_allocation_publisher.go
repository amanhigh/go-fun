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

	extras := map[string]string{
		fun.MetadataEnrollmentID: enrollment.ID,
		fun.MetadataPersonID:     enrollment.PersonID,
	}

	return sap.base.PublishWithExtras(ctx, fun.TopicAllocateSeatCmd, payload, extras)
}

func (sap *seatAllocationPublisher) SeatReserved(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	payload := fun.SeatReservedEvtV1{
		EnrollmentID: enrollment.ID,
		PersonID:     enrollment.PersonID,
		Grade:        enrollment.Grade,
		ReservedAt:   time.Now().UTC(),
	}

	extras := map[string]string{
		fun.MetadataEnrollmentID: enrollment.ID,
		fun.MetadataPersonID:     enrollment.PersonID,
	}

	return sap.base.PublishWithExtras(ctx, fun.TopicSeatReservedEvt, payload, extras)
}

func (sap *seatAllocationPublisher) SeatWaitlisted(ctx context.Context, enrollment fun.Enrollment, reason string) common.HttpError {
	payload := fun.SeatWaitlistedEvtV1{
		EnrollmentID: enrollment.ID,
		PersonID:     enrollment.PersonID,
		Grade:        enrollment.Grade,
		Reason:       reason,
		WaitlistedAt: time.Now().UTC(),
	}

	extras := map[string]string{
		fun.MetadataEnrollmentID: enrollment.ID,
		fun.MetadataPersonID:     enrollment.PersonID,
	}

	return sap.base.PublishWithExtras(ctx, fun.TopicSeatWaitlistedEvt, payload, extras)
}
