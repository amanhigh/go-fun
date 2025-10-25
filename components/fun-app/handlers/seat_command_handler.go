package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/fun"
)

// SeatCommandHandler handles seat-related commands and events.
type SeatCommandHandler interface {
	AllocateSeatCmd(msg *message.Message) error
	SeatReservedEvt(msg *message.Message) error
	SeatWaitlistedEvt(msg *message.Message) error
}

type SeatCommandHandlerImpl struct {
	SeatManager       manager.SeatManagerInterface       `container:"type"`
	EnrollmentManager manager.EnrollmentManagerInterface `container:"type"`
}

func NewSeatCommandHandler() *SeatCommandHandlerImpl { return &SeatCommandHandlerImpl{} }

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
	return h.EnrollmentManager.ConfirmFlow(ctx, e)
}

func (h *SeatCommandHandlerImpl) SeatWaitlistedEvt(msg *message.Message) error {
	var evt fun.SeatWaitlistedEvtV1
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return fmt.Errorf("unmarshal seat waitlisted evt: %w", err)
	}
	ctx := stampCtx(msg.Context(), msg.Metadata, evt.EnrollmentID, msg.UUID)
	return h.SeatManager.OnSeatWaitlistedEvt(ctx, evt)
}

// emit helpers removed; direct publisher calls are used.
