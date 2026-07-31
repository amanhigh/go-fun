package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/fun"
)

const defaultSeatAllocationFailureReason = "unknown allocation failure"

// SeatMessageHandler handles seat-related commands and events.
type SeatMessageHandler interface {
	HandleAllocateSeatCmd(msg *message.Message) error
	HandleSeatReservedEvt(msg *message.Message) error
	HandleSeatWaitlistedEvt(msg *message.Message) error
	HandleSeatAllocationFailedEvt(msg *message.Message) error
	HandleDeadLetteredAllocateSeatCmd(msg *message.Message) error
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
	cmd, err := util.DecodeAndValidateMessage[fun.AllocateSeatCmdV1](msg)
	if err != nil {
		return err
	}

	return h.SeatManager.AllocateSeat(msg.Context(), cmd)
}

func (h *SeatMessageHandlerImpl) HandleSeatReservedEvt(msg *message.Message) error {
	evt, err := util.DecodeAndValidateMessage[fun.SeatReservedEvtV1](msg)
	if err != nil {
		return err
	}
	e := fun.Enrollment{ID: evt.EnrollmentID, PersonID: evt.PersonID, Grade: evt.Grade}
	return h.EnrollmentManager.OnSeatReservedEvt(msg.Context(), e)
}

func (h *SeatMessageHandlerImpl) HandleSeatWaitlistedEvt(msg *message.Message) error {
	evt, err := util.DecodeAndValidateMessage[fun.SeatWaitlistedEvtV1](msg)
	if err != nil {
		return err
	}
	enrollment := fun.Enrollment{ID: evt.EnrollmentID, PersonID: evt.PersonID, Grade: evt.Grade}
	// FIXME: Waitlisted flow ends here — add exponential backoff retry to re-check seat availability before terminal failure.
	// Persist WAITLISTED state via manager sink (idempotent).
	return h.EnrollmentManager.UpdateToWaitlisted(msg.Context(), enrollment)
}

// HandleSeatAllocationFailedEvt compensates a failed seat allocation by cancelling the enrollment.
func (h *SeatMessageHandlerImpl) HandleSeatAllocationFailedEvt(msg *message.Message) error {
	evt, err := util.DecodeAndValidateMessage[fun.SeatAllocationFailedEvtV1](msg)
	if err != nil {
		return err
	}
	return h.EnrollmentManager.CancelEnrollment(msg.Context(), evt.EnrollmentID)
}

// HandleDeadLetteredAllocateSeatCmd publishes the allocation failure event.
func (h *SeatMessageHandlerImpl) HandleDeadLetteredAllocateSeatCmd(msg *message.Message) error {
	var cmd fun.AllocateSeatCmdV1
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		return fmt.Errorf("unmarshal dead-lettered allocate seat cmd v1: %w", err)
	}
	reason := ""
	if msg.Metadata != nil {
		reason = msg.Metadata.Get(middleware.ReasonForPoisonedKey)
	}
	if reason == "" {
		reason = defaultSeatAllocationFailureReason
	}
	return h.SeatManager.PublishSeatAllocationFailed(msg.Context(), fun.Enrollment{ID: cmd.EnrollmentID, PersonID: cmd.PersonID}, reason)
}
