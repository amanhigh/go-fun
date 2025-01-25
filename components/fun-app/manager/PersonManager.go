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

type PersonManagerInterface interface {
	CreatePerson(c context.Context, request fun.PersonRequest) (person fun.Person, err common.HttpError)
	DeletePerson(c context.Context, id string) (err common.HttpError)
	UpdatePerson(c context.Context, id string, request fun.PersonRequest) (err common.HttpError)

	ListPersons(c context.Context, query fun.PersonQuery) (response fun.PersonList, err common.HttpError)
	GetPerson(c context.Context, id string) (person fun.Person, err common.HttpError)
	ListPersonAudit(c context.Context, id string) (response []fun.PersonAudit, err common.HttpError)
}

type PersonManager struct {
	Dao    dao.PersonDaoInterface
	Tracer trace.Tracer
}

func NewPersonManager(dao dao.PersonDaoInterface, tracer trace.Tracer) PersonManagerInterface {
	return &PersonManager{Dao: dao, Tracer: tracer}
}

// CreatePerson creates a new person in the PersonManager.
//
// It takes two parameters:
// - c: a context.Context object representing the current context.
// - person: Person object representing the person to be created.
//
// It returns two values:
// - id: a string representing the ID of the newly created person.
// - err: an error representing any error that occurred during the creation process.
func (p *PersonManager) CreatePerson(c context.Context, request fun.PersonRequest) (person fun.Person, err common.HttpError) {
	subLogger := zerolog.Ctx(c).With().Str("Name", request.Name).Int("Age", request.Age).Str("Gender", request.Gender).Logger()

	ctx, span := p.Tracer.Start(c, "CreatePerson.Manager")
	defer span.End()

	/* Create Person */
	person.Name = request.Name
	person.Age = request.Age
	person.Gender = request.Gender

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err common.HttpError) {
		if err = p.Dao.Create(c, &person); err == nil {
			subLogger.Info().Ctx(c).Str("Id", person.Id).Msg("Person Created")
		}
		return
	})

	return
}

func (p *PersonManager) ListPersons(c context.Context, personQuery fun.PersonQuery) (response fun.PersonList, err common.HttpError) {
	ctx, span := p.Tracer.Start(c, "ListPersons.Manager", trace.WithAttributes(
		attribute.String("gender", personQuery.Gender),
		attribute.String("name", personQuery.Name),
		attribute.Int("offset", personQuery.Offset),
		attribute.Int("limit", personQuery.Limit),
	))
	defer span.End()

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err common.HttpError) {
		response, err = p.Dao.ListPerson(c, personQuery)
		return
	})
	return
}

func (p *PersonManager) ListPersonAudit(c context.Context, id string) (response []fun.PersonAudit, err common.HttpError) {
	ctx, span := p.Tracer.Start(c, "GetPersonAudit.Manager", trace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err1 common.HttpError) {
		response, err1 = p.Dao.ListPersonAudit(c, id)
		return
	})
	return
}

func (p *PersonManager) GetPerson(c context.Context, id string) (person fun.Person, err common.HttpError) {
	ctx, span := p.Tracer.Start(c, "GetPerson.Manager", trace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err common.HttpError) {
		return p.Dao.FindById(c, id, &person)
	})
	return
}

func (p *PersonManager) UpdatePerson(c context.Context, id string, request fun.PersonRequest) (err common.HttpError) {
	//Create Person
	var person fun.Person
	person.Id = id
	person.Name = request.Name
	person.Age = request.Age
	person.Gender = request.Gender

	ctx, span := p.Tracer.Start(c, "UpdatePerson.Manager", trace.WithAttributes(
		attribute.String("id", id),
		attribute.String("Name", request.Name),
		attribute.Int("Age", request.Age),
		attribute.String("Gender", request.Gender),
	))
	defer span.End()

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err common.HttpError) {
		err = p.Dao.Update(c, &person)
		return
	})
	return
}

func (p *PersonManager) DeletePerson(c context.Context, id string) (err common.HttpError) {
	var person fun.Person
	ctx, span := p.Tracer.Start(c, "DeletePerson.Manager", trace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	err = p.Dao.UseOrCreateTx(ctx, func(c context.Context) (err common.HttpError) {
		if person, err = p.GetPerson(c, id); err == nil {
			span.AddEvent("Person Found for Deletion", trace.WithAttributes(attribute.String("Name", person.Name))) //Adds message in log section of Span
			/* Delete from DB */
			err = p.Dao.DeleteById(c, id, &person)
		}
		return
	})

	return
}
