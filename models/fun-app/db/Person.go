package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Person struct {
	// Validations - https://gin-gonic.com/docs/examples/binding-and-validation/
	Id      string `gorm:"primaryKey"`
	Name    string `gorm:"not null" binding:"required,min=1,max=25,name=person"`
	Age     int    `gorm:"not null" binding:"required,min=1,max=150"`
	Gender  string `gorm:"not null" binding:"required,eq=MALE|eq=FEMALE" enums:"MALE,FEMALE"`
	Version int64  `gorm:"not null" json:"-" binding:"-"`
}

func (p *Person) BeforeCreate(tx *gorm.DB) (err error) {
	p.Id = uuid.NewString()[:8]
	return
}
