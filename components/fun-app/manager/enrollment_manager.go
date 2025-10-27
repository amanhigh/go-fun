package manager

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

// EnrollmentManagerInterface orchestrates enrollment flows and delegates seat allocation.
//
// Architecture rules enforced here:
//   - Flow is always Handler -> Manager -> Publisher. Handlers never publish directly.
//   - Managers only talk to their own publisher; cross-domain messages use Manager-to-Manager calls.
//   - EnrollmentManager responsibilities:
//   - EnrollPerson persists an initiated enrollment and publishes EnrollCmd (C1).
//   - EnrollCmd delegates seat allocation to SeatManager, which publishes AllocateSeat (C2).
//   - OnSeatReservedEvt, UpdateToWaitlisted, OnEnrollmentConfirmedEvt, and OnEnrollmentCancelledEvt are idempotent sinks that persist status without publishing.
//   - SeatManager publishes only seat-related commands/events and never touches enrollment publishers.
//
// TODO: Rename Person usage to Student once the domain model is updated.
type EnrollmentManagerInterface interface {
	EnrollPerson(ctx context.Context, personID string, grade int) (fun.Enrollment, common.HttpError)
	GetEnrollment(ctx context.Context, personID string) (fun.Enrollment, common.HttpError)
	EnrollCmd(ctx context.Context, cmd fun.EnrollCmdV1, meta message.Metadata, messageID string) common.HttpError
	OnSeatReservedEvt(ctx context.Context, enrollment fun.Enrollment) common.HttpError
	UpdateToWaitlisted(ctx context.Context, enrollment fun.Enrollment) common.HttpError
	OnEnrollmentConfirmedEvt(ctx context.Context, evt fun.EnrollmentConfirmedEvtV1) common.HttpError
	OnEnrollmentCancelledEvt(ctx context.Context, evt fun.EnrollmentCancelledEvtV1) common.HttpError
}

type EnrollmentManager struct {
	PersonManager       PersonManagerInterface
	EnrollmentDao       dao.EnrollmentDaoInterface
	EnrollmentPublisher publisher.EnrollmentPublisher
	SeatManager         SeatManagerInterface
}

func NewEnrollmentManager(
	personManager PersonManagerInterface,
	enrollmentDao dao.EnrollmentDaoInterface,
	enrollmentPublisher publisher.EnrollmentPublisher,
	seatManager SeatManagerInterface,
) EnrollmentManagerInterface {
	return &EnrollmentManager{
		PersonManager:       personManager,
		EnrollmentDao:       enrollmentDao,
		EnrollmentPublisher: enrollmentPublisher,
		SeatManager:         seatManager,
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

// EnrollCmd coordinates seat allocation by delegating to SeatManager.
func (em *EnrollmentManager) EnrollCmd(ctx context.Context, cmd fun.EnrollCmdV1, meta message.Metadata, messageID string) common.HttpError {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = em.bindMessageContext(ctx, cmd.EnrollmentID, meta, messageID)

	enrollment := fun.Enrollment{
		ID:       cmd.EnrollmentID,
		PersonID: cmd.PersonID,
		Grade:    cmd.Grade,
	}

	return em.SeatManager.PublishAllocateSeat(ctx, enrollment)
}

// OnSeatReservedEvt persists CONFIRMED status without publishing.
func (em *EnrollmentManager) OnSeatReservedEvt(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	return em.updateStatusByID(ctx, enrollment.ID, fun.EnrollmentStatusConfirmed)
}

// UpdateToWaitlisted persists WAITLISTED status without publishing.
func (em *EnrollmentManager) UpdateToWaitlisted(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	return em.updateStatusByID(ctx, enrollment.ID, fun.EnrollmentStatusWaitlisted)
}

// OnEnrollmentConfirmedEvt persists CONFIRMED status without publishing.
func (em *EnrollmentManager) OnEnrollmentConfirmedEvt(ctx context.Context, evt fun.EnrollmentConfirmedEvtV1) common.HttpError {
	return em.updateStatusByID(ctx, evt.EnrollmentID, fun.EnrollmentStatusConfirmed)
}

// OnEnrollmentCancelledEvt persists CANCELLED status without publishing.
func (em *EnrollmentManager) OnEnrollmentCancelledEvt(ctx context.Context, evt fun.EnrollmentCancelledEvtV1) common.HttpError {
	return em.updateStatusByID(ctx, evt.EnrollmentID, fun.EnrollmentStatusCancelled)
}

func (em *EnrollmentManager) bindMessageContext(ctx context.Context, enrollmentID string, meta message.Metadata, messageID string) context.Context {
	correlationID := enrollmentID
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

func (em *EnrollmentManager) updateStatusByID(ctx context.Context, enrollmentID, status string) common.HttpError {
	return em.EnrollmentDao.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		var enrollment fun.Enrollment
		if findErr := em.EnrollmentDao.FindById(c, enrollmentID, &enrollment); findErr != nil {
			return findErr
		}
		if enrollment.Status == status {
			return nil
		}
		enrollment.Status = status
		return em.EnrollmentDao.Update(c, &enrollment)
	})
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
	return &fun.Enrollment{
		PersonID: personID,
		Grade:    grade,
		Status:   fun.EnrollmentStatusSeatAllocationInitiated,
	}
}
