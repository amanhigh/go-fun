package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type PersonManagerInterface interface {
	CreatePerson(c context.Context, person fun.PersonRequest) (id string, err common.HttpError)
	DeletePerson(c context.Context, id string) (err common.HttpError)
	UpdatePerson(c context.Context, id string, request fun.PersonRequest) (err common.HttpError)

	ListPersons(c context.Context, personQuery fun.PersonQuery) (response fun.PersonList, err common.HttpError)
	GetPerson(c context.Context, id string) (person fun.Person, err common.HttpError)
}

type PersonManager struct {
	Dao    dao.PersonDaoInterface `inject:""`
	Tracer trace.Tracer           `inject:""`
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
func (self *PersonManager) CreatePerson(c context.Context, personRequest fun.PersonRequest) (id string, err common.HttpError) {
	personFields := log.Fields{"Name": personRequest.Name, "Age": personRequest.Age, "Gender": personRequest.Gender}

	/* Create Person */
	var person fun.Person
	person.Name = personRequest.Name
	person.Age = personRequest.Age
	person.Gender = personRequest.Gender

	err = self.Dao.UseOrCreateTx(c, func(c context.Context) (err common.HttpError) {
		if err = self.Dao.Create(c, &person); err == nil {
			id = person.Id
			log.WithContext(c).WithField("Id", id).WithFields(personFields).Info("Person Created")
		}
		return
	})

	return
}

func (self *PersonManager) ListPersons(c context.Context, personQuery fun.PersonQuery) (response fun.PersonList, err common.HttpError) {
	err = self.Dao.UseOrCreateTx(c, func(c context.Context) (err common.HttpError) {
		response, err = self.Dao.ListPerson(c, personQuery)
		return
	})
	return
}

func (self *PersonManager) GetPerson(c context.Context, id string) (person fun.Person, err common.HttpError) {
	ctx, span := self.Tracer.Start(c, "GetPerson.Manager", trace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	err = self.Dao.UseOrCreateTx(ctx, func(c context.Context) (err common.HttpError) {
		return self.Dao.FindById(c, id, &person)
	})
	return
}

func (self *PersonManager) UpdatePerson(c context.Context, id string, request fun.PersonRequest) (err common.HttpError) {
	//Create Person
	var person fun.Person
	person.Id = id
	person.Name = request.Name
	person.Age = request.Age
	person.Gender = request.Gender

	err = self.Dao.UseOrCreateTx(c, func(c context.Context) (err common.HttpError) {
		err = self.Dao.Update(c, &person)
		return
	})
	return
}

func (self *PersonManager) DeletePerson(c context.Context, id string) (err common.HttpError) {
	var person fun.Person

	err = self.Dao.UseOrCreateTx(c, func(c context.Context) (err common.HttpError) {
		if person, err = self.GetPerson(c, id); err == nil {
			/* Delete from DB */
			err = self.Dao.DeleteById(c, id, &person)
		}
		return
	})

	return
}
