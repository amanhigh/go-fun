package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/fun"
)

// SeatCommandHandler handles seat-related commands and events.
type SeatCommandHandler interface {
	AllocateSeatCmd(msg *message.Message) error
	SeatReservedEvt(msg *message.Message) error
	SeatWaitlistedEvt(msg *message.Message) error
	PoisonAllocate(msg *message.Message) error
}

type SeatCommandHandlerImpl struct {
	SeatManager       manager.SeatManagerInterface
	EnrollmentManager manager.EnrollmentManagerInterface
}

// NewSeatCommandHandler constructs handler with explicit dependencies.
func NewSeatCommandHandler(seatManager manager.SeatManagerInterface, enrollmentManager manager.EnrollmentManagerInterface) *SeatCommandHandlerImpl {
	return &SeatCommandHandlerImpl{SeatManager: seatManager, EnrollmentManager: enrollmentManager}
}

var _ SeatCommandHandler = (*SeatCommandHandlerImpl)(nil)

func (h *SeatCommandHandlerImpl) AllocateSeatCmd(msg *message.Message) error {
	var cmd fun.AllocateSeatCmdV1
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		return fmt.Errorf("unmarshal allocate seat cmd: %w", err)
	}

	ctx := stampCtx(msg.Context(), msg.Metadata, cmd.EnrollmentID, msg.UUID)
	return h.SeatManager.AllocateSeat(ctx, cmd)
}

func (h *SeatCommandHandlerImpl) SeatReservedEvt(msg *message.Message) error {
	var evt fun.SeatReservedEvtV1
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return fmt.Errorf("unmarshal seat reserved evt: %w", err)
	}
	ctx := stampCtx(msg.Context(), msg.Metadata, evt.EnrollmentID, msg.UUID)
	e := fun.Enrollment{ID: evt.EnrollmentID, PersonID: evt.PersonID, Grade: evt.Grade}
	return h.EnrollmentManager.OnSeatReservedEvt(ctx, e)
}

func (h *SeatCommandHandlerImpl) SeatWaitlistedEvt(msg *message.Message) error {
	var evt fun.SeatWaitlistedEvtV1
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return fmt.Errorf("unmarshal seat waitlisted evt: %w", err)
	}
	ctx := stampCtx(msg.Context(), msg.Metadata, evt.EnrollmentID, msg.UUID)
	enrollment := fun.Enrollment{ID: evt.EnrollmentID, PersonID: evt.PersonID, Grade: evt.Grade}
	// 1) Persist WAITLISTED state via manager sink (idempotent)
	if err := h.EnrollmentManager.UpdateToWaitlisted(ctx, enrollment); err != nil {
		return err
	}
	// BUG: Not Sure if Retry would work Unclear Who does SeatAllocation Retry.
	// Done after sink; retry will be driven by middleware on AllocateSeatCmd path
	return nil
}

// PoisonAllocate consumes poison messages after retries and cancels enrollment.
func (h *SeatCommandHandlerImpl) PoisonAllocate(msg *message.Message) error {
	var cmd fun.AllocateSeatCmdV1
	// HACK: Rename PoisonAllocateCmd ?
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		return fmt.Errorf("unmarshal poison allocate: %w", err)
	}
	ctx := stampCtx(msg.Context(), msg.Metadata, cmd.EnrollmentID, msg.UUID)
	return h.EnrollmentManager.CancelEnrollmentAndPublish(ctx, fun.EnrollmentCancelledEvtV1{
		EnrollmentID: cmd.EnrollmentID,
		PersonID:     cmd.PersonID,
		Reason:       "waitlist_retries_exhausted",
		CancelledAt:  time.Now().UTC(),
	})
}

// emit helpers removed; direct publisher calls are used.
