package fun

import (
	"github.com/amanhigh/go-fun/models/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PersonRequest struct {
	// Validations - https://gin-gonic.com/docs/examples/binding-and-validation/
	Name   string `json:"name" gorm:"not null" binding:"required,min=1,max=25,name=person"`
	Age    int    `json:"age" gorm:"not null" binding:"required,min=1,max=150"`
	Gender string `json:"gender" gorm:"not null" binding:"required,eq=MALE|eq=FEMALE" enums:"MALE,FEMALE"`
}

type PersonPath struct {
	Id string `uri:"id" binding:"required"`
}

type PersonQuery struct {
	common.Pagination
	Name   string `form:"name" binding:"omitempty,min=1,max=25,name=person"`
	Gender string `form:"gender" binding:"omitempty,eq=MALE|eq=FEMALE"`
}

type PersonList struct {
	Records []Person

	common.PaginatedResponse
}

type Person struct {
	PersonRequest
	Id      string `gorm:"primaryKey" json:"id"`
	Version int64  `gorm:"not null" json:"-" binding:"-"`
}

func (p *Person) BeforeCreate(tx *gorm.DB) (err error) {
	p.Id = uuid.NewString()[:8]
	return
}
