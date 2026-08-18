package manager

import (
	"context"

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
//   - EnrollStudent persists an initiated enrollment and publishes EnrollCmd (C1).
//   - EnrollCmd delegates seat allocation to SeatManager, which publishes AllocateSeat (C2).
//   - OnSeatReservedEvt persists CONFIRMED status.
//   - CancelEnrollment persists CANCELLED status for allocation failures.
//   - SeatManager publishes only seat-related commands/events and never touches enrollment publishers.
//
// TODO: Rename Student usage to Student once the domain model is updated.
type EnrollmentManagerInterface interface {
	EnrollStudent(ctx context.Context, studentID string, grade int) (fun.Enrollment, common.HttpError)
	GetEnrollment(ctx context.Context, studentID string) (fun.Enrollment, common.HttpError)
	EnrollCmd(ctx context.Context, cmd fun.EnrollCmdV1) common.HttpError
	OnSeatReservedEvt(ctx context.Context, enrollment fun.Enrollment) common.HttpError
	UpdateToWaitlisted(ctx context.Context, enrollment fun.Enrollment) common.HttpError
	CancelEnrollment(ctx context.Context, enrollmentID string) common.HttpError
}

type EnrollmentManager struct {
	StudentManager       StudentManagerInterface
	EnrollmentDao       dao.EnrollmentDaoInterface
	EnrollmentPublisher publisher.EnrollmentPublisher
	SeatManager         SeatManagerInterface
}

func NewEnrollmentManager(
	studentManager StudentManagerInterface,
	enrollmentDao dao.EnrollmentDaoInterface,
	enrollmentPublisher publisher.EnrollmentPublisher,
	seatManager SeatManagerInterface,
) *EnrollmentManager {
	return &EnrollmentManager{
		StudentManager:       studentManager,
		EnrollmentDao:       enrollmentDao,
		EnrollmentPublisher: enrollmentPublisher,
		SeatManager:         seatManager,
	}
}

var _ EnrollmentManagerInterface = (*EnrollmentManager)(nil)

func (em *EnrollmentManager) EnrollStudent(ctx context.Context, studentID string, grade int) (fun.Enrollment, common.HttpError) {
	student, err := em.StudentManager.GetStudent(ctx, studentID)
	if err != nil {
		return fun.Enrollment{}, err
	}

	enrollment := em.buildEnrollment(student.Id, grade)
	if err := em.upsertEnrollment(ctx, enrollment); err != nil {
		return fun.Enrollment{}, err
	}

	if publishErr := em.EnrollmentPublisher.Enroll(ctx, *enrollment); publishErr != nil {
		return fun.Enrollment{}, publishErr
	}
	return *enrollment, nil
}

func (em *EnrollmentManager) GetEnrollment(ctx context.Context, studentID string) (fun.Enrollment, common.HttpError) {
	var enrollment fun.Enrollment
	if err := em.EnrollmentDao.FindByStudentID(ctx, studentID, &enrollment); err != nil {
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
		StudentID: cmd.StudentID,
		Grade:    cmd.Grade,
	}

	return em.SeatManager.PublishAllocateSeat(ctx, enrollment)
}

// OnSeatReservedEvt persists CONFIRMED status without publishing a redundant event.
func (em *EnrollmentManager) OnSeatReservedEvt(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	return em.updateStatusByID(ctx, enrollment.ID, fun.EnrollmentStatusConfirmed)
}

// UpdateToWaitlisted persists WAITLISTED status without publishing.
func (em *EnrollmentManager) UpdateToWaitlisted(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	return em.updateStatusByID(ctx, enrollment.ID, fun.EnrollmentStatusWaitlisted)
}

// CancelEnrollment persists CANCELLED status without publishing a cancellation event.
func (em *EnrollmentManager) CancelEnrollment(ctx context.Context, enrollmentID string) common.HttpError {
	return em.updateStatusByID(ctx, enrollmentID, fun.EnrollmentStatusCancelled)
}

func (em *EnrollmentManager) updateStatusByID(ctx context.Context, enrollmentID, status string) common.HttpError {
	var persisted fun.Enrollment

	err := em.EnrollmentDao.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		if findErr := em.EnrollmentDao.FindById(c, enrollmentID, &persisted); findErr != nil {
			return findErr
		}
		if persisted.Status == status {
			return nil
		}
		if persisted.Status != fun.EnrollmentStatusSeatAllocationInitiated {
			return nil
		}
		persisted.Status = status
		return em.EnrollmentDao.Update(c, &persisted)
	})
	return err
}

func (em *EnrollmentManager) upsertEnrollment(ctx context.Context, enrollment *fun.Enrollment) common.HttpError {
	return em.EnrollmentDao.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		var existing fun.Enrollment
		err := em.EnrollmentDao.FindByStudentID(c, enrollment.StudentID, &existing)
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

func (em *EnrollmentManager) buildEnrollment(studentID string, grade int) *fun.Enrollment {
	return &fun.Enrollment{
		StudentID: studentID,
		Grade:    grade,
		Status:   fun.EnrollmentStatusSeatAllocationInitiated,
	}
}
