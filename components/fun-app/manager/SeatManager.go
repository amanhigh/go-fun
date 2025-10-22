package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog"
)

// SeatManagerInterface encapsulates seat allocation transitions for the enrollment saga.
//
// AllocateSeat applies seat allocation rules for an enrollment and emits the appropriate
// events. It is idempotent and safe to call multiple times for the same enrollment.
type SeatManagerInterface interface {
	AllocateSeat(ctx context.Context, enrollment fun.Enrollment) (fun.Enrollment, common.HttpError)
}

// SeatManager is the concrete implementation coordinating state changes and events.
type SeatManager struct {
	SeatPublisher publisher.SeatAllocationPublisher
}

const seatWaitlistThreshold = 5 // TODO: move to config when real capacity is implemented

// NewSeatManager constructs a SeatManager.
func NewSeatManager(seatPub publisher.SeatAllocationPublisher) SeatManagerInterface {
	return &SeatManager{
		SeatPublisher: seatPub,
	}
}

// AllocateSeat evaluates capacity and advances the enrollment through seat-related states.
// Draft implementation keeps the current threshold rule: grades < 5 are available; otherwise waitlist.
func (sm *SeatManager) AllocateSeat(ctx context.Context, enrollment fun.Enrollment) (fun.Enrollment, common.HttpError) {

	// Stateless decision based on grade threshold; emit corresponding events.
	fresh := enrollment
	// No-op if already terminal.
	if fresh.Status == fun.EnrollmentStatusConfirmed || fresh.Status == fun.EnrollmentStatusWaitlisted {
		return fresh, nil
	}

	if fresh.Grade >= seatWaitlistThreshold {
		fresh.Status = fun.EnrollmentStatusWaitlisted
		if err := sm.SeatPublisher.SeatWaitlisted(ctx, fresh, "capacity_unavailable"); err != nil {
			return fun.Enrollment{}, err
		}
		zerolog.Ctx(ctx).Info().Str("enrollmentId", fresh.ID).Int("grade", fresh.Grade).Msg("Enrollment waitlisted by seat manager")
		return fresh, nil
	}

	// Seat available path: reserve then confirm.
	reserved := fresh
	reserved.Status = fun.EnrollmentStatusSeatReserved
	if err := sm.SeatPublisher.SeatReserved(ctx, reserved); err != nil {
		return fun.Enrollment{}, err
	}

	zerolog.Ctx(ctx).Info().Str("enrollmentId", reserved.ID).Int("grade", reserved.Grade).Msg("Seat reserved by seat manager")
	return reserved, nil
}
