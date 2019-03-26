package manager

import (
	log "github.com/Sirupsen/logrus"
	"github.com/amanhigh/go-fun/apps/models/fun-app/db"
	"github.com/jinzhu/gorm"
)

type PersonManagerInterface interface {
	CreatePerson(person db.Person) (err error)
	DeletePerson(id string) (err error)

	GetAllPersons() (persons []db.Person, err error)
}

type PersonManager struct {
	Db *gorm.DB `inject:""`
}

func (self *PersonManager) CreatePerson(person db.Person) (err error) {
	personFields := log.Fields{"Name": person.Name}

	/*
		Create new Person
	*/
	if err = self.Db.Create(&person).Error; err != nil {
		log.WithFields(personFields).WithField("Error", err).Error("Error Creating Person")
	}
	return
}

func (self *PersonManager) GetAllPersons() (persons []db.Person, err error) {
	err = self.Db.Find(&persons, &db.Person{}).Error
	return
}

func (self *PersonManager) DeletePerson(id string) (err error) {
	var person = db.Person{}

	/* Find Person in DB */
	if err = self.Db.Find(&person, "person_id=?", id).Error; err == nil {

		/* Delete from DB */
		if err == nil {
			err = self.Db.Delete(&person).Error
		}
	}

	return
}
