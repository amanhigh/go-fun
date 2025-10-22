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

// EnrollmentCommandHandler handles enrollment-saga commands such as Enroll.
type EnrollmentCommandHandler interface {
	Handle(msg *message.Message) error
}

type enrollmentCommandHandler struct {
	mgr manager.EnrollmentManagerInterface
}

// NewEnrollmentCommandHandler constructs the handler.
func NewEnrollmentCommandHandler(mgr manager.EnrollmentManagerInterface) EnrollmentCommandHandler {
	return &enrollmentCommandHandler{mgr: mgr}
}

// Handle processes EnrollCmdV1 messages and delegates to the manager.
func (h *enrollmentCommandHandler) Handle(msg *message.Message) error {
	var payload fun.EnrollCmdV1
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		zerolog.Ctx(msg.Context()).Error().
			Err(err).
			Str("handler", "enrollment_enroll_requested").
			Msg("Failed to unmarshal EnrollCmdV1 payload")
		return fmt.Errorf("unmarshal enroll cmd: %w", err)
	}

	handlerCtx := msg.Context()
	if handlerCtx == nil {
		handlerCtx = context.Background()
	}

	corr := payload.EnrollmentID
	if v := msg.Metadata.Get(common.MetadataCorrelationIDKey); v != "" {
		corr = v
	}
	handlerCtx = common.WithCorrelation(handlerCtx, corr)
	if caus := msg.Metadata.Get(common.MetadataCausationIDKey); caus != "" {
		handlerCtx = common.WithCausation(handlerCtx, caus)
	} else if msg.UUID != "" {
		handlerCtx = common.WithCausation(handlerCtx, msg.UUID)
	}

	return h.mgr.ProcessEnrollRequested(handlerCtx, payload, msg.Metadata, msg.UUID)
}
