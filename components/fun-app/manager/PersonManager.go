package manager

import (
	"context"

	db2 "github.com/amanhigh/go-fun/models/fun-app/db"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PersonManagerInterface interface {
	CreatePerson(c context.Context, person db2.Person) (err error)
	DeletePerson(c context.Context, id string) (err error)

	GetAllPersons(c context.Context) (persons []db2.Person, err error)
	GetPerson(c context.Context, name string) (person db2.Person, err error)
}

type PersonManager struct {
	Db *gorm.DB `inject:""`
}

func (self *PersonManager) CreatePerson(c context.Context, person db2.Person) (err error) {
	personFields := log.Fields{"Name": person.Name}

	/*
		Create new Person
	*/
	if err = self.Db.Create(&person).Error; err != nil {
		log.WithContext(c).WithFields(personFields).WithField("Error", err).Error("Error Creating Person")
	} else {
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

	/* Find Person in DB */
	if err = self.Db.Find(&person, "person_id=?", id).Error; err == nil {

		/* Delete from DB */
		if err == nil {
			err = self.Db.Delete(&person).Error
		}
	}

	return
}
