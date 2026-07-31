package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/fun"
)

// EnrollmentMessageHandler handles enrollment saga commands/events.
type EnrollmentMessageHandler interface {
	HandleEnrollCmd(msg *message.Message) error
	HandleEnrollmentConfirmedEvt(msg *message.Message) error
	HandleEnrollmentCancelledEvt(msg *message.Message) error
}

type EnrollmentMessageHandlerImpl struct {
	Manager manager.EnrollmentManagerInterface
}

// NewEnrollmentMessageHandler constructs handler with explicit manager dependency.
func NewEnrollmentMessageHandler(manager manager.EnrollmentManagerInterface) *EnrollmentMessageHandlerImpl {
	return &EnrollmentMessageHandlerImpl{Manager: manager}
}

var _ EnrollmentMessageHandler = (*EnrollmentMessageHandlerImpl)(nil)

// HandleEnrollCmd forwards EnrollCmdV1 to EnrollmentManager; it delegates to SeatManager internally.
func (h *EnrollmentMessageHandlerImpl) HandleEnrollCmd(msg *message.Message) error {
	cmd, err := util.DecodeAndValidateMessage[fun.EnrollCmdV1](msg)
	if err != nil {
		return err
	}

	return h.Manager.EnrollCmd(msg.Context(), cmd)
}

// HandleEnrollmentConfirmedEvt persists CONFIRMED status via manager sink.
func (h *EnrollmentMessageHandlerImpl) HandleEnrollmentConfirmedEvt(msg *message.Message) error {
	var evt fun.EnrollmentConfirmedEvtV1
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return fmt.Errorf("unmarshal enrollment confirmed evt: %w", err)
	}
	return h.Manager.OnEnrollmentConfirmedEvt(msg.Context(), evt)
}

// HandleEnrollmentCancelledEvt persists CANCELLED status via manager sink.
func (h *EnrollmentMessageHandlerImpl) HandleEnrollmentCancelledEvt(msg *message.Message) error {
	var evt fun.EnrollmentCancelledEvtV1
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return fmt.Errorf("unmarshal enrollment cancelled evt: %w", err)
	}
	return h.Manager.OnEnrollmentCancelledEvt(msg.Context(), evt)
}
