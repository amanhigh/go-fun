package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/fun"
)

// SeatMessageHandler handles seat-related commands and events.
type SeatMessageHandler interface {
	HandleAllocateSeatCmd(msg *message.Message) error
	HandleSeatReservedEvt(msg *message.Message) error
	HandleSeatWaitlistedEvt(msg *message.Message) error
	HandlePoisonedAllocateSeatCmd(msg *message.Message) error
}

type SeatMessageHandlerImpl struct {
	SeatManager       manager.SeatManagerInterface
	EnrollmentManager manager.EnrollmentManagerInterface
}

// NewSeatMessageHandler constructs handler with explicit dependencies.
func NewSeatMessageHandler(seatManager manager.SeatManagerInterface, enrollmentManager manager.EnrollmentManagerInterface) *SeatMessageHandlerImpl {
	return &SeatMessageHandlerImpl{SeatManager: seatManager, EnrollmentManager: enrollmentManager}
}

var _ SeatMessageHandler = (*SeatMessageHandlerImpl)(nil)

func (h *SeatMessageHandlerImpl) HandleAllocateSeatCmd(msg *message.Message) error {
	var cmd fun.AllocateSeatCmdV1
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		return fmt.Errorf("unmarshal allocate seat cmd: %w", err)
	}

	ctx := stampCtx(msg.Context(), msg.Metadata, cmd.EnrollmentID, msg.UUID)
	return h.SeatManager.AllocateSeat(ctx, cmd)
}

func (h *SeatMessageHandlerImpl) HandleSeatReservedEvt(msg *message.Message) error {
	var evt fun.SeatReservedEvtV1
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return fmt.Errorf("unmarshal seat reserved evt: %w", err)
	}
	ctx := stampCtx(msg.Context(), msg.Metadata, evt.EnrollmentID, msg.UUID)
	e := fun.Enrollment{ID: evt.EnrollmentID, StudentID: evt.StudentID, Grade: evt.Grade}
	return h.EnrollmentManager.OnSeatReservedEvt(ctx, e)
}

func (h *SeatMessageHandlerImpl) HandleSeatWaitlistedEvt(msg *message.Message) error {
	var evt fun.SeatWaitlistedEvtV1
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return fmt.Errorf("unmarshal seat waitlisted evt: %w", err)
	}
	ctx := stampCtx(msg.Context(), msg.Metadata, evt.EnrollmentID, msg.UUID)
	enrollment := fun.Enrollment{ID: evt.EnrollmentID, StudentID: evt.StudentID, Grade: evt.Grade}
	// Persist WAITLISTED state via manager sink (idempotent).
	return h.EnrollmentManager.UpdateToWaitlisted(ctx, enrollment)
}

// HandlePoisonedAllocateSeatCmd handles a poisoned AllocateSeatCmdV1 by cancelling the enrollment.
func (h *SeatMessageHandlerImpl) HandlePoisonedAllocateSeatCmd(msg *message.Message) error {
	var cmd fun.AllocateSeatCmdV1
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		return fmt.Errorf("unmarshal poisoned allocate seat cmd v1: %w", err)
	}
	ctx := stampCtx(msg.Context(), msg.Metadata, cmd.EnrollmentID, msg.UUID)
	return h.EnrollmentManager.CancelEnrollmentAndPublish(ctx, fun.EnrollmentCancelledEvtV1{
		EnrollmentID: cmd.EnrollmentID,
		StudentID:     cmd.StudentID,
		Reason:       fun.EnrollmentCancellationReasonSeatAllocationFailed,
		CancelledAt:  time.Now().UTC(),
	})
}

// emit helpers removed; direct publisher calls are used.
