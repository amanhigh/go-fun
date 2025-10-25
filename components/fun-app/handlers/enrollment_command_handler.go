package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/components/fun-app/publisher"
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
	Manager             manager.EnrollmentManagerInterface `container:"type"`
	SeatPublisher       publisher.SeatAllocationPublisher  `container:"type"`
	EnrollmentPublisher publisher.EnrollmentPublisher      `container:"type"`
	EnrollmentDao       dao.EnrollmentDaoInterface         `container:"type"`
}

func NewEnrollmentCommandHandler() *EnrollmentCommandHandlerImpl {
	return &EnrollmentCommandHandlerImpl{}
}

// EnrollCmd processes EnrollCmdV1 commands via EnrollmentManager, then emits follow-ups.
func (h *EnrollmentCommandHandlerImpl) EnrollCmd(msg *message.Message) error {
	var payload fun.EnrollCmdV1
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("unmarshal enroll cmd: %w", err)
	}

	ctx := stampCtx(msg.Context(), msg.Metadata, payload.EnrollmentID, msg.UUID)

	actions, err := h.Manager.ProcessEnrollRequested(ctx, payload, msg.Metadata, msg.UUID)
	if err != nil {
		return err
	}

	if actions.Waitlisted {
		enrollment := fun.Enrollment{ID: payload.EnrollmentID, PersonID: payload.PersonID, Grade: payload.Grade, Status: fun.EnrollmentStatusWaitlisted}
		return h.SeatPublisher.SeatWaitlisted(ctx, enrollment, "capacity_unavailable")
	}
	if actions.AllocationStarted {
		enrollment := fun.Enrollment{ID: payload.EnrollmentID, PersonID: payload.PersonID, Grade: payload.Grade, Status: fun.EnrollmentStatusSeatAllocationInitiated}
		return h.SeatPublisher.AllocateSeat(ctx, enrollment)
	}
	return nil
}

// EnrollmentConfirmedEvt idempotently persists CONFIRMED status.
func (h *EnrollmentCommandHandlerImpl) EnrollmentConfirmedEvt(msg *message.Message) error {
	var evt fun.EnrollmentConfirmedEvtV1
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return fmt.Errorf("unmarshal enrollment confirmed evt: %w", err)
	}
	ctx := stampCtx(msg.Context(), msg.Metadata, evt.EnrollmentID, msg.UUID)

	return h.EnrollmentDao.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		var enrollment fun.Enrollment
		if findErr := h.EnrollmentDao.FindById(c, evt.EnrollmentID, &enrollment); findErr != nil {
			return findErr
		}
		if enrollment.Status == fun.EnrollmentStatusConfirmed {
			return nil
		}
		enrollment.Status = fun.EnrollmentStatusConfirmed
		return h.EnrollmentDao.Update(c, &enrollment)
	})
}

// EnrollmentCancelledEvt persists CANCELLED status.
func (h *EnrollmentCommandHandlerImpl) EnrollmentCancelledEvt(msg *message.Message) error {
	var evt fun.EnrollmentCancelledEvtV1
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return fmt.Errorf("unmarshal enrollment cancelled evt: %w", err)
	}
	ctx := stampCtx(msg.Context(), msg.Metadata, evt.EnrollmentID, msg.UUID)

	return h.EnrollmentDao.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		var enrollment fun.Enrollment
		if findErr := h.EnrollmentDao.FindById(c, evt.EnrollmentID, &enrollment); findErr != nil {
			return findErr
		}
		enrollment.Status = fun.EnrollmentStatusCancelled
		return h.EnrollmentDao.Update(c, &enrollment)
	})
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
