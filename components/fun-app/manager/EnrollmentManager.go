package manager

import (
	"context"
	"net/http"

	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// EnrollmentManagerInterface orchestrates enrollment flows using existing person records.
// TODO: Rename Person usage to Student once the domain model is updated.
type EnrollmentManagerInterface interface {
	EnrollPerson(ctx context.Context, personID string, grade int) (fun.Enrollment, common.HttpError)
	GetEnrollment(ctx context.Context, personID string) (fun.EnrollmentResponse, common.HttpError)
}

type EnrollmentManager struct {
	PersonManager PersonManagerInterface
	EnrollmentDao dao.EnrollmentDaoInterface
	Tracer        trace.Tracer
}

func NewEnrollmentManager(personManager PersonManagerInterface, enrollmentDao dao.EnrollmentDaoInterface, tracer trace.Tracer) EnrollmentManagerInterface {
	return &EnrollmentManager{
		PersonManager: personManager,
		EnrollmentDao: enrollmentDao,
		Tracer:        tracer,
	}
}

func (em *EnrollmentManager) EnrollPerson(ctx context.Context, personID string, grade int) (fun.Enrollment, common.HttpError) {
	ctx, span := em.Tracer.Start(ctx, "EnrollPerson.Manager")
	defer span.End()

	// FIXME: Move Grade to Person/Student
	person, err := em.PersonManager.GetPerson(ctx, personID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fun.Enrollment{}, err
	}

	if grade >= 5 {
		conflictErr := common.NewHttpError("SeatUnavailable", http.StatusConflict)
		span.RecordError(conflictErr)
		span.SetStatus(codes.Error, conflictErr.Error())
		zerolog.Ctx(ctx).Info().Str("personId", personID).Int("grade", grade).Msg("Enrollment conflict: seat unavailable")
		return fun.Enrollment{}, conflictErr
	}

	enrollment := &fun.Enrollment{
		PersonID: person.Id,
		Grade:    grade,
		Status:   fun.EnrollmentStatusActive,
	}

	if err = em.upsertEnrollment(ctx, enrollment); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fun.Enrollment{}, err
	}

	span.SetStatus(codes.Ok, "Enrollment completed")
	zerolog.Ctx(ctx).Info().Str("personId", personID).Int("grade", grade).Msg("Enrollment completed")
	return *enrollment, nil
}

func (em *EnrollmentManager) GetEnrollment(ctx context.Context, personID string) (fun.EnrollmentResponse, common.HttpError) {
	var enrollment fun.Enrollment
	if err := em.EnrollmentDao.FindByPersonID(ctx, personID, &enrollment); err != nil {
		return fun.EnrollmentResponse{}, err
	}

	return fun.EnrollmentResponse{
		EnrollmentID: enrollment.ID,
		PersonID:     enrollment.PersonID,
		Grade:        enrollment.Grade,
		Status:       enrollment.Status,
	}, nil
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
