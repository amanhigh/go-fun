package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

// EnrollmentCommandHandler handles enrollment saga commands/events.
type EnrollmentCommandHandler interface {
	EnrollCmd(msg *message.Message) error
	EnrollmentConfirmedEvt(msg *message.Message) error
	EnrollmentCancelledEvt(msg *message.Message) error
}

type EnrollmentCommandHandlerImpl struct {
	Manager manager.EnrollmentManagerInterface
}

// NewEnrollmentCommandHandler constructs handler with explicit manager dependency and returns interface.
func NewEnrollmentCommandHandler(manager manager.EnrollmentManagerInterface) EnrollmentCommandHandler {
	return &EnrollmentCommandHandlerImpl{Manager: manager}
}

// EnrollCmd forwards EnrollCmdV1 to EnrollmentManager; it delegates to SeatManager internally.
func (h *EnrollmentCommandHandlerImpl) EnrollCmd(msg *message.Message) error {
	var cmd fun.EnrollCmdV1
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		return fmt.Errorf("unmarshal enroll cmd: %w", err)
	}

	ctx := stampCtx(msg.Context(), msg.Metadata, cmd.EnrollmentID, msg.UUID)
	return h.Manager.EnrollCmd(ctx, cmd, msg.Metadata, msg.UUID)
}

// EnrollmentConfirmedEvt persists CONFIRMED status via manager sink.
func (h *EnrollmentCommandHandlerImpl) EnrollmentConfirmedEvt(msg *message.Message) error {
	var evt fun.EnrollmentConfirmedEvtV1
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return fmt.Errorf("unmarshal enrollment confirmed evt: %w", err)
	}
	ctx := stampCtx(msg.Context(), msg.Metadata, evt.EnrollmentID, msg.UUID)
	return h.Manager.OnEnrollmentConfirmedEvt(ctx, evt)
}

// EnrollmentCancelledEvt persists CANCELLED status via manager sink.
func (h *EnrollmentCommandHandlerImpl) EnrollmentCancelledEvt(msg *message.Message) error {
	var evt fun.EnrollmentCancelledEvtV1
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return fmt.Errorf("unmarshal enrollment cancelled evt: %w", err)
	}
	ctx := stampCtx(msg.Context(), msg.Metadata, evt.EnrollmentID, msg.UUID)
	return h.Manager.OnEnrollmentCancelledEvt(ctx, evt)
}

// stampCtx helper to apply correlation/causation from message metadata.
func stampCtx(in context.Context, meta message.Metadata, enrollmentID, messageID string) context.Context {
	if in == nil {
		in = context.Background()
	}
	corr := enrollmentID
	if meta != nil {
		if v := meta.Get(common.MetadataCorrelationIDKey); v != "" {
			corr = v
		}
	}
	out := common.WithCorrelation(in, corr)
	causation := messageID
	if meta != nil {
		if v := meta.Get(common.MetadataCausationIDKey); v != "" {
			causation = v
		}
	}
	if causation != "" {
		out = common.WithCausation(out, causation)
	}
	return out
}
