package manager

import (
	"context"
	"net/http"
	"sync"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// EnrollmentManagerInterface orchestrates enrollment flows using existing person records.
// TODO: Rename Person usage to Student once the domain model is updated.
type EnrollmentManagerInterface interface {
	EnrollPerson(ctx context.Context, personID string, grade int) (fun.Person, common.HttpError)
	GetEnrollment(ctx context.Context, personID string) (fun.EnrollmentResponse, common.HttpError)
}

type EnrollmentManager struct {
	PersonManager PersonManagerInterface
	Tracer        trace.Tracer
	grades        map[string]int
	mu            sync.RWMutex
}

func NewEnrollmentManager(personManager PersonManagerInterface, tracer trace.Tracer) EnrollmentManagerInterface {
	return &EnrollmentManager{
		PersonManager: personManager,
		Tracer:        tracer,
		grades:        make(map[string]int),
	}
}

func (em *EnrollmentManager) EnrollPerson(ctx context.Context, personID string, grade int) (fun.Person, common.HttpError) {
	ctx, span := em.Tracer.Start(ctx, "EnrollPerson.Manager")
	defer span.End()

	// FIXME: Move Grade to Person/Student
	person, err := em.PersonManager.GetPerson(ctx, personID)
	if err != nil {
		span.RecordError(err)
		return fun.Person{}, err
	}

	if grade >= 5 {
		conflictErr := common.NewHttpError("SeatUnavailable", http.StatusConflict)
		span.RecordError(conflictErr)
		zerolog.Ctx(ctx).Info().Str("personId", personID).Int("grade", grade).Msg("Enrollment conflict: seat unavailable")
		return fun.Person{}, conflictErr
	}

	em.mu.Lock()
	em.grades[person.Id] = grade
	em.mu.Unlock()

	zerolog.Ctx(ctx).Info().Str("personId", personID).Int("grade", grade).Msg("Enrollment completed")
	return person, nil
}

func (em *EnrollmentManager) GetEnrollment(ctx context.Context, personID string) (fun.EnrollmentResponse, common.HttpError) {
	em.mu.RLock()
	grade, ok := em.grades[personID]
	em.mu.RUnlock()
	if !ok {
		notFound := common.NewHttpError("EnrollmentNotFound", http.StatusNotFound)
		return fun.EnrollmentResponse{}, notFound
	}

	person, err := em.PersonManager.GetPerson(ctx, personID)
	if err != nil {
		return fun.EnrollmentResponse{}, err
	}

	return fun.EnrollmentResponse{PersonID: person.Id, Grade: grade, Status: "ACTIVE"}, nil
}
