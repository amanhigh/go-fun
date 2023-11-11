package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/models/common"
	db2 "github.com/amanhigh/go-fun/models/fun-app/db"
	"github.com/amanhigh/go-fun/models/fun-app/server"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PersonManagerInterface interface {
	CreatePerson(c context.Context, person server.PersonRequest) (id string, err common.HttpError)
	DeletePerson(c context.Context, id string) (err common.HttpError)

	GetAllPersons(c context.Context) (persons []db2.Person, err error)
	GetPerson(c context.Context, id string) (person db2.Person, err common.HttpError)
}

type PersonManager struct {
	Db  *gorm.DB               `inject:""`
	Dao dao.PersonDaoInterface `inject:""`
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
func (self *PersonManager) CreatePerson(c context.Context, person server.PersonRequest) (id string, err common.HttpError) {
	personFields := log.Fields{"Name": person.Person.Name, "Age": person.Person.Age, "Gender": person.Person.Gender}

	err = self.Dao.UseOrCreateTx(c, func(c context.Context) (err common.HttpError) {
		if err = self.Dao.Create(c, &person.Person); err == nil {
			id = person.Id
			log.WithContext(c).WithFields(personFields).Info("Person Created")
		}
		return
	})

	return
}

func (self *PersonManager) GetAllPersons(c context.Context) (persons []db2.Person, err error) {
	err = self.Db.Find(&persons).Error
	return
}

func (self *PersonManager) GetPerson(c context.Context, id string) (person db2.Person, err common.HttpError) {
	err = self.Dao.UseOrCreateTx(c, func(c context.Context) (err common.HttpError) {
		return self.Dao.FindById(c, id, &person)
	})
	return
}

func (self *PersonManager) DeletePerson(c context.Context, id string) (err common.HttpError) {
	var person db2.Person

	err = self.Dao.UseOrCreateTx(c, func(c context.Context) (err common.HttpError) {
		if person, err = self.GetPerson(c, id); err == nil {
			/* Delete from DB */
			err = self.Dao.DeleteById(c, id, &person)
		}
		return
	})

	return
}
