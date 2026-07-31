package handlers

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/fun"
)

// EnrollmentMessageHandler handles enrollment commands only.
type EnrollmentMessageHandler interface {
	HandleEnrollCmd(msg *message.Message) error
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
