package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

// FIXME: #B Add to AGENTS.md
// EnrollmentManagerInterface orchestrates enrollment flows and delegates seat allocation.
//
// Architecture rules enforced here:
//   - Flow is always Handler -> Manager -> Publisher. Handlers never publish directly.
//   - Managers only talk to their own publisher; cross-domain messages use Manager-to-Manager calls.
//   - EnrollmentManager responsibilities:
//   - EnrollPerson persists an initiated enrollment and publishes EnrollCmd (C1).
//   - EnrollCmd delegates seat allocation to SeatManager, which publishes AllocateSeat (C2).
//   - OnSeatReservedEvt persists status and emits EnrollmentConfirmedEvt.
//   - CancelEnrollmentAndPublish persists status and emits EnrollmentCancelledEvt for origin-side cancellations.
//   - OnEnrollmentConfirmedEvt and OnEnrollmentCancelledEvt are idempotent sinks that persist status without publishing.
//   - SeatManager publishes only seat-related commands/events and never touches enrollment publishers.
//
// TODO: #C Rename Person usage to Student once the domain model is updated.
type EnrollmentManagerInterface interface {
	EnrollPerson(ctx context.Context, personID string, grade int) (fun.Enrollment, common.HttpError)
	GetEnrollment(ctx context.Context, personID string) (fun.Enrollment, common.HttpError)
	EnrollCmd(ctx context.Context, cmd fun.EnrollCmdV1) common.HttpError
	OnSeatReservedEvt(ctx context.Context, enrollment fun.Enrollment) common.HttpError
	UpdateToWaitlisted(ctx context.Context, enrollment fun.Enrollment) common.HttpError
	CancelEnrollmentAndPublish(ctx context.Context, evt fun.EnrollmentCancelledEvtV1) common.HttpError
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
) *EnrollmentManager {
	return &EnrollmentManager{
		PersonManager:       personManager,
		EnrollmentDao:       enrollmentDao,
		EnrollmentPublisher: enrollmentPublisher,
		SeatManager:         seatManager,
	}
}

var _ EnrollmentManagerInterface = (*EnrollmentManager)(nil)

func (em *EnrollmentManager) EnrollPerson(ctx context.Context, personID string, grade int) (fun.Enrollment, common.HttpError) {
	person, err := em.PersonManager.GetPerson(ctx, personID)
	if err != nil {
		return fun.Enrollment{}, err
	}

	enrollment := em.buildEnrollment(person.Id, grade)
	if err := em.upsertEnrollment(ctx, enrollment); err != nil {
		return fun.Enrollment{}, err
	}

	if common.CorrelationFrom(ctx) == "" {
		ctx = common.WithCorrelation(ctx, enrollment.ID)
	}

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
func (em *EnrollmentManager) EnrollCmd(ctx context.Context, cmd fun.EnrollCmdV1) common.HttpError {
	if ctx == nil {
		ctx = context.Background()
	}

	enrollment := fun.Enrollment{
		ID:       cmd.EnrollmentID,
		PersonID: cmd.PersonID,
		Grade:    cmd.Grade,
	}

	return em.SeatManager.PublishAllocateSeat(ctx, enrollment)
}

// OnSeatReservedEvt persists CONFIRMED status and publishes confirmation when status changes.
func (em *EnrollmentManager) OnSeatReservedEvt(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	persisted, changed, err := em.updateStatusByID(ctx, enrollment.ID, fun.EnrollmentStatusConfirmed)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	return em.EnrollmentPublisher.EnrollmentConfirmedEvt(ctx, persisted)
}

// UpdateToWaitlisted persists WAITLISTED status without publishing.
func (em *EnrollmentManager) UpdateToWaitlisted(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	_, _, err := em.updateStatusByID(ctx, enrollment.ID, fun.EnrollmentStatusWaitlisted)
	return err
}

// CancelEnrollmentAndPublish persists CANCELLED status and emits cancellation event when status changes.
func (em *EnrollmentManager) CancelEnrollmentAndPublish(ctx context.Context, evt fun.EnrollmentCancelledEvtV1) common.HttpError {
	persisted, changed, err := em.updateStatusByID(ctx, evt.EnrollmentID, fun.EnrollmentStatusCancelled)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	return em.EnrollmentPublisher.EnrollmentCancelledEvt(ctx, persisted, evt.Reason)
}

// OnEnrollmentConfirmedEvt persists CONFIRMED status without publishing.
func (em *EnrollmentManager) OnEnrollmentConfirmedEvt(ctx context.Context, evt fun.EnrollmentConfirmedEvtV1) common.HttpError {
	_, _, err := em.updateStatusByID(ctx, evt.EnrollmentID, fun.EnrollmentStatusConfirmed)
	return err
}

// OnEnrollmentCancelledEvt persists CANCELLED status without publishing.
func (em *EnrollmentManager) OnEnrollmentCancelledEvt(ctx context.Context, evt fun.EnrollmentCancelledEvtV1) common.HttpError {
	_, _, err := em.updateStatusByID(ctx, evt.EnrollmentID, fun.EnrollmentStatusCancelled)
	return err
}

func (em *EnrollmentManager) updateStatusByID(ctx context.Context, enrollmentID, status string) (fun.Enrollment, bool, common.HttpError) {
	var persisted fun.Enrollment
	changed := false

	err := em.EnrollmentDao.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		if findErr := em.EnrollmentDao.FindById(c, enrollmentID, &persisted); findErr != nil {
			return findErr
		}
		if persisted.Status == status {
			return nil
		}
		persisted.Status = status
		changed = true
		return em.EnrollmentDao.Update(c, &persisted)
	})
	if err != nil {
		return fun.Enrollment{}, false, err
	}
	return persisted, changed, nil
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

func (em *EnrollmentManager) buildEnrollment(personID string, grade int) *fun.Enrollment {
	return &fun.Enrollment{
		PersonID: personID,
		Grade:    grade,
		Status:   fun.EnrollmentStatusSeatAllocationInitiated,
	}
}
