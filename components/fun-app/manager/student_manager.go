package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type StudentManagerInterface interface {
	CreateStudent(c context.Context, request fun.StudentRequest) (student fun.Student, err common.HttpError)
	DeleteStudent(c context.Context, id string) (err common.HttpError)
	UpdateStudent(c context.Context, id string, request fun.StudentRequest) (err common.HttpError)

	ListStudents(c context.Context, query fun.StudentQuery) (response fun.StudentList, err common.HttpError)
	GetStudent(c context.Context, id string) (student fun.Student, err common.HttpError)
	ListStudentAudit(c context.Context, id string) (response []fun.StudentAudit, err common.HttpError)
}

type StudentManager struct {
	Dao    dao.StudentDaoInterface
	Tracer trace.Tracer
}

func NewStudentManager(dao dao.StudentDaoInterface, tracer trace.Tracer) *StudentManager {
	return &StudentManager{Dao: dao, Tracer: tracer}
}

var _ StudentManagerInterface = (*StudentManager)(nil)

// CreateStudent creates a new student in the StudentManager.
//
// It takes two parameters:
// - c: a context.Context object representing the current context.
// - student: Student object representing the student to be created.
//
// It returns two values:
// - id: a string representing the ID of the newly created student.
// - err: an error representing any error that occurred during the creation process.
func (p *StudentManager) CreateStudent(c context.Context, request fun.StudentRequest) (student fun.Student, err common.HttpError) {
	subLogger := zerolog.Ctx(c).With().Str("Name", request.Name).Int("Age", request.Age).Str("Gender", request.Gender).Logger()

	ctx, span := p.Tracer.Start(c, "CreateStudent.Manager")
	defer span.End()

	/* Create Student */
	student.Name = request.Name
	student.Age = request.Age
	student.Gender = request.Gender

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err common.HttpError) {
		if err = p.Dao.Create(c, &student); err == nil {
			subLogger.Info().Ctx(c).Str("Id", student.Id).Msg("Student Created")
		}
		return
	})

	return
}

func (p *StudentManager) ListStudents(c context.Context, studentQuery fun.StudentQuery) (response fun.StudentList, err common.HttpError) {
	ctx, span := p.Tracer.Start(c, "ListStudents.Manager", trace.WithAttributes(
		attribute.String("gender", studentQuery.Gender),
		attribute.String("name", studentQuery.Name),
		attribute.Int("offset", studentQuery.Offset),
		attribute.Int("limit", studentQuery.Limit),
	))
	defer span.End()

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err common.HttpError) {
		response, err = p.Dao.ListStudent(c, studentQuery)
		return
	})
	return
}

func (p *StudentManager) ListStudentAudit(c context.Context, id string) (response []fun.StudentAudit, err common.HttpError) {
	ctx, span := p.Tracer.Start(c, "GetStudentAudit.Manager", trace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err1 common.HttpError) {
		response, err1 = p.Dao.ListStudentAudit(c, id)
		return
	})
	return
}

func (p *StudentManager) GetStudent(c context.Context, id string) (student fun.Student, err common.HttpError) {
	ctx, span := p.Tracer.Start(c, "GetStudent.Manager", trace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err common.HttpError) {
		return p.Dao.FindById(c, id, &student)
	})
	return
}

func (p *StudentManager) UpdateStudent(c context.Context, id string, request fun.StudentRequest) (err common.HttpError) {
	// Create Student
	var student fun.Student
	student.Id = id
	student.Name = request.Name
	student.Age = request.Age
	student.Gender = request.Gender

	ctx, span := p.Tracer.Start(c, "UpdateStudent.Manager", trace.WithAttributes(
		attribute.String("id", id),
		attribute.String("Name", request.Name),
		attribute.Int("Age", request.Age),
		attribute.String("Gender", request.Gender),
	))
	defer span.End()

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err common.HttpError) {
		err = p.Dao.Update(c, &student)
		return
	})
	return
}

func (p *StudentManager) DeleteStudent(c context.Context, id string) (err common.HttpError) {
	ctx, span := p.Tracer.Start(c, "DeleteStudent.Manager", trace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err common.HttpError) {
		span.AddEvent("Student Found for Deletion", trace.WithAttributes(attribute.String("id", id)))
		// BUG: Deleting with an empty fun.Student{} prevents AfterDelete audit association with the deleted student.
		return p.Dao.DeleteById(c, id, &fun.Student{})
	})

	return
}
