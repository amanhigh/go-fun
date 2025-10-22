package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog"
)

// SeatCommandHandler defines the contract for handling seat-related commands.
//
// Implementations must be idempotent and safe for at-least-once delivery.
type SeatCommandHandler interface {
	// Handle processes an incoming AllocateSeat command message.
	Handle(msg *message.Message) error
}

type seatCommandHandler struct {
	seatMgr manager.SeatManagerInterface
}

// NewSeatCommandHandler constructs a SeatCommandHandler.
func NewSeatCommandHandler(seatMgr manager.SeatManagerInterface) SeatCommandHandler {
	return &seatCommandHandler{seatMgr: seatMgr}
}

// Handle processes AllocateSeatCmdV1 messages and invokes SeatManager.
func (h *seatCommandHandler) Handle(msg *message.Message) error {
	var cmd fun.AllocateSeatCmdV1
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		zerolog.Ctx(msg.Context()).Error().
			Err(err).
			Str("handler", "seat_allocate").
			Msg("Failed to unmarshal AllocateSeatCmdV1 payload")
		return fmt.Errorf("unmarshal allocate seat cmd: %w", err)
	}

	handlerCtx := msg.Context()
	if handlerCtx == nil {
		handlerCtx = context.Background()
	}

	// Correlation and causation propagation from metadata/message.
	corr := cmd.EnrollmentID
	if v := msg.Metadata.Get(common.MetadataCorrelationIDKey); v != "" {
		corr = v
	}
	handlerCtx = common.WithCorrelation(handlerCtx, corr)
	if caus := msg.Metadata.Get(common.MetadataCausationIDKey); caus != "" {
		handlerCtx = common.WithCausation(handlerCtx, caus)
	} else if msg.UUID != "" {
		handlerCtx = common.WithCausation(handlerCtx, msg.UUID)
	}

	// Build minimal enrollment DTO for seat manager.
	enrollment := fun.Enrollment{
		ID:       cmd.EnrollmentID,
		PersonID: cmd.PersonID,
		Grade:    cmd.Grade,
		Status:   fun.EnrollmentStatusSeatAllocationInitiated,
	}

	if _, err := h.seatMgr.AllocateSeat(handlerCtx, enrollment); err != nil {
		return err
	}
	return nil
}
