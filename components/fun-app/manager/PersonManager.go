package manager

import (
	"context"

	db2 "github.com/amanhigh/go-fun/models/fun-app/db"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PersonManagerInterface interface {
	CreatePerson(c context.Context, person db2.Person) (id string, err error)
	DeletePerson(c context.Context, id string) (err error)

	GetAllPersons(c context.Context) (persons []db2.Person, err error)
	GetPerson(c context.Context, name string) (person db2.Person, err error)
}

type PersonManager struct {
	Db *gorm.DB `inject:""`
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
func (self *PersonManager) CreatePerson(c context.Context, person db2.Person) (id string, err error) {
	personFields := log.Fields{"Name": person.Name}

	if err = self.Db.Create(&person).Error; err != nil {
		log.WithContext(c).WithFields(personFields).WithField("Error", err).Error("Error Creating Person")
	} else {
		id = person.Id
		log.WithContext(c).WithFields(personFields).Info("Person Created")
	}
	return
}

func (self *PersonManager) GetAllPersons(c context.Context) (persons []db2.Person, err error) {
	err = self.Db.Find(&persons).Error
	return
}

func (self *PersonManager) GetPerson(c context.Context, name string) (person db2.Person, err error) {
	err = self.Db.First(&person, "name = ?", name).Error
	return
}

func (self *PersonManager) DeletePerson(c context.Context, id string) (err error) {
	var person = db2.Person{}

	//TODO: Fix Delete.
	/* Find Person in DB */
	if err = self.Db.Find(&person, "person_id=?", id).Error; err == nil {

		/* Delete from DB */
		if err == nil {
			err = self.Db.Delete(&person).Error
		}
	}

	return
}
