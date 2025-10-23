package manager

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog"
)

// EnrollmentManagerInterface orchestrates enrollment flows using existing person records.
// TODO: Rename Person usage to Student once the domain model is updated.
type EnrollmentManagerInterface interface {
	EnrollPerson(ctx context.Context, personID string, grade int) (fun.Enrollment, common.HttpError)
	GetEnrollment(ctx context.Context, personID string) (fun.Enrollment, common.HttpError)
	ProcessEnrollRequested(ctx context.Context, event fun.EnrollCmdV1, meta message.Metadata, messageID string) (EnrollmentActions, common.HttpError)
	// Seat flows (called by SeatManager)
	SeatReservedFlow(ctx context.Context, e fun.Enrollment) common.HttpError
	SeatWaitlistedFlow(ctx context.Context, e fun.Enrollment, reason string) common.HttpError
	ConfirmFlow(ctx context.Context, e fun.Enrollment) common.HttpError
	// Update only, used when consuming SeatWaitlisted event (no publish)
	UpdateToWaitlisted(ctx context.Context, e fun.Enrollment) common.HttpError
}

type EnrollmentManager struct {
	PersonManager       PersonManagerInterface
	EnrollmentDao       dao.EnrollmentDaoInterface
	EnrollmentPublisher publisher.EnrollmentPublisher
	SeatPublisher       publisher.SeatAllocationPublisher
}

const (
	waitlistGradeThreshold       = 5
	seatWaitlistedReasonCapacity = "capacity_unavailable"
)

// transitionState captures actions to emit after DB state changes.
type transitionState struct {
	emitAllocationStarted bool
	emitWaitlisted        bool
	noOp                  bool
}

// EnrollmentActions distills transition intents for handlers.
type EnrollmentActions struct {
	AllocationStarted bool
	Waitlisted        bool
}

func NewEnrollmentManager(
	personManager PersonManagerInterface,
	enrollmentDao dao.EnrollmentDaoInterface,
	enrollmentPublisher publisher.EnrollmentPublisher,
	seatPublisher publisher.SeatAllocationPublisher,
) EnrollmentManagerInterface {
	return &EnrollmentManager{
		PersonManager:       personManager,
		EnrollmentDao:       enrollmentDao,
		EnrollmentPublisher: enrollmentPublisher,
		SeatPublisher:       seatPublisher,
	}
}

func (em *EnrollmentManager) EnrollPerson(ctx context.Context, personID string, grade int) (fun.Enrollment, common.HttpError) {
	person, err := em.retrievePerson(ctx, personID)
	if err != nil {
		return fun.Enrollment{}, err
	}

	enrollment := em.buildEnrollment(person.Id, grade)
	if err := em.upsertEnrollment(ctx, enrollment); err != nil {
		return fun.Enrollment{}, err
	}

	ctx = common.WithCorrelation(ctx, enrollment.ID)

	if publishErr := em.EnrollmentPublisher.Enroll(ctx, *enrollment); publishErr != nil {
		return fun.Enrollment{}, publishErr
	}
	return *enrollment, nil
}

func (em *EnrollmentManager) GetEnrollment(ctx context.Context, personID string) (fun.Enrollment, common.HttpError) {
	var enrollment fun.Enrollment
	if err := em.EnrollmentDao.FindByPersonID(ctx, personID, &enrollment); err != nil {
		return fun.Enrollment{}, err
	}

	return enrollment, nil
}

func (em *EnrollmentManager) ProcessEnrollRequested(ctx context.Context, event fun.EnrollCmdV1, meta message.Metadata, messageID string) (EnrollmentActions, common.HttpError) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = em.bindMessageContext(ctx, event, meta, messageID)

	state, _, err := em.computeTransitions(ctx, event)
	if err != nil {
		if state.noOp {
			return EnrollmentActions{}, nil
		}
		return EnrollmentActions{}, err
	}
	if state.noOp {
		return EnrollmentActions{}, nil
	}

	actions := EnrollmentActions{
		AllocationStarted: state.emitAllocationStarted,
		Waitlisted:        state.emitWaitlisted,
	}
	// Do not publish here; handlers will emit commands/events based on updated state.
	return actions, nil
}

func (em *EnrollmentManager) bindMessageContext(ctx context.Context, event fun.EnrollCmdV1, meta message.Metadata, messageID string) context.Context {
	correlationID := event.EnrollmentID
	if meta != nil {
		if corr := meta.Get(common.MetadataCorrelationIDKey); corr != "" {
			correlationID = corr
		}
	}

	causationID := messageID
	if meta != nil {
		if causation := meta.Get(common.MetadataCausationIDKey); causation != "" {
			causationID = causation
		}
	}

	ctx = common.WithCorrelation(ctx, correlationID)
	if causationID != "" {
		ctx = common.WithCausation(ctx, causationID)
	}
	return ctx
}

func (em *EnrollmentManager) evalAndUpdateTransition(ctx context.Context, enrollment *fun.Enrollment) (transitionState, common.HttpError) {
	var state transitionState

	switch enrollment.Status {
	case fun.EnrollmentStatusConfirmed, fun.EnrollmentStatusWaitlisted:
		state.noOp = true
		return state, nil
	}

	if enrollment.Grade >= waitlistGradeThreshold {
		if err := em.updateEnrollmentStatus(ctx, enrollment, fun.EnrollmentStatusWaitlisted); err != nil {
			return state, err
		}
		state.emitWaitlisted = true
		return state, nil
	}

	if enrollment.Status == fun.EnrollmentStatusSeatAllocationInitiated {
		// No DB change; seat allocation will be processed asynchronously.
		state.emitAllocationStarted = true
	}

	return state, nil
}

func (em *EnrollmentManager) computeTransitions(ctx context.Context, event fun.EnrollCmdV1) (transitionState, fun.Enrollment, common.HttpError) {
	var (
		state      transitionState
		enrollment fun.Enrollment
	)

	err := em.EnrollmentDao.UseOrCreateTx(ctx, func(txCtx context.Context) common.HttpError {
		if err := em.EnrollmentDao.FindById(txCtx, event.EnrollmentID, &enrollment); err != nil {
			return err
		}

		var evalErr common.HttpError
		state, evalErr = em.evalAndUpdateTransition(txCtx, &enrollment)
		return evalErr
	})

	return state, enrollment, err
}

// No publishing here; handlers own emitting commands/events.

// SeatReservedFlow publishes SeatReserved event; persistence to Confirmed happens on reserved evt.
func (em *EnrollmentManager) SeatReservedFlow(ctx context.Context, e fun.Enrollment) common.HttpError {
	return em.SeatPublisher.SeatReserved(ctx, e)
}

// SeatWaitlistedFlow persists WAITLISTED and publishes SeatWaitlisted.
func (em *EnrollmentManager) SeatWaitlistedFlow(ctx context.Context, e fun.Enrollment, reason string) common.HttpError {
	if err := em.UpdateToWaitlisted(ctx, e); err != nil {
		return err
	}
	return em.SeatPublisher.SeatWaitlisted(ctx, e, reason)
}

// UpdateToWaitlisted persists WAITLISTED without publishing (used by SeatManager on sink event).
func (em *EnrollmentManager) UpdateToWaitlisted(ctx context.Context, e fun.Enrollment) common.HttpError {
	return em.EnrollmentDao.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		var existing fun.Enrollment
		if findErr := em.EnrollmentDao.FindById(c, e.ID, &existing); findErr != nil {
			return findErr
		}
		if existing.Status != fun.EnrollmentStatusWaitlisted {
			existing.Status = fun.EnrollmentStatusWaitlisted
			if updErr := em.EnrollmentDao.Update(c, &existing); updErr != nil {
				return updErr
			}
		}
		e = existing
		return nil
	})
}

// ConfirmFlow persists CONFIRMED and publishes EnrollmentConfirmed.
func (em *EnrollmentManager) ConfirmFlow(ctx context.Context, e fun.Enrollment) common.HttpError {
	if err := em.EnrollmentDao.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		var existing fun.Enrollment
		if findErr := em.EnrollmentDao.FindById(c, e.ID, &existing); findErr != nil {
			return findErr
		}
		if existing.Status != fun.EnrollmentStatusConfirmed {
			existing.Status = fun.EnrollmentStatusConfirmed
			if updErr := em.EnrollmentDao.Update(c, &existing); updErr != nil {
				return updErr
			}
		}
		e = existing
		return nil
	}); err != nil {
		return err
	}
	return em.EnrollmentPublisher.EnrollmentConfirmedEvt(ctx, e)
}

func (em *EnrollmentManager) upsertEnrollment(ctx context.Context, enrollment *fun.Enrollment) common.HttpError {
	return em.EnrollmentDao.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		var existing fun.Enrollment
		err := em.EnrollmentDao.FindByPersonID(c, enrollment.PersonID, &existing)
		switch err {
		case nil:
			existing.Grade = enrollment.Grade
			existing.Status = enrollment.Status
			updateErr := em.EnrollmentDao.Update(c, &existing)
			if updateErr == nil {
				*enrollment = existing
			}
			return updateErr
		case common.ErrNotFound:
			return em.EnrollmentDao.Create(c, enrollment)
		default:
			return err
		}
	})
}

func (em *EnrollmentManager) retrievePerson(ctx context.Context, personID string) (fun.Person, common.HttpError) {
	return em.PersonManager.GetPerson(ctx, personID)
}

func (em *EnrollmentManager) buildEnrollment(personID string, grade int) *fun.Enrollment {
	status := fun.EnrollmentStatusSeatAllocationInitiated
	if grade >= waitlistGradeThreshold {
		status = fun.EnrollmentStatusWaitlisted
	}
	return &fun.Enrollment{
		PersonID: personID,
		Grade:    grade,
		Status:   status,
	}
}

func (em *EnrollmentManager) updateEnrollmentStatus(ctx context.Context, enrollment *fun.Enrollment, status string) common.HttpError {
	if enrollment.Status == status {
		return nil
	}

	enrollment.Status = status
	if err := em.EnrollmentDao.Update(ctx, enrollment); err != nil {
		return err
	}

	em.logStatusTransition(ctx, enrollment)
	return nil
}

func (em *EnrollmentManager) logStatusTransition(ctx context.Context, enrollment *fun.Enrollment) {
	zerolog.Ctx(ctx).Info().
		Str("enrollmentId", enrollment.ID).
		Str("personId", enrollment.PersonID).
		Str("status", enrollment.Status).
		Msg("Enrollment status updated")
}
