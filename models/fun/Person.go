package fun

import (
	"time"

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
	common.Sort
	Name   string `form:"name" binding:"omitempty,min=1,max=25,name=person"`
	Gender string `form:"gender" binding:"omitempty,eq=MALE|eq=FEMALE"`
}

type PersonList struct {
	Records  []Person                 `json:"records"`
	Metadata common.PaginatedResponse `json:"metadata"`
}

type Person struct {
	PersonRequest
	Id string `gorm:"primaryKey" json:"id"`
}

func (p *Person) BeforeCreate(tx *gorm.DB) (err error) {
	p.Id = uuid.NewString()[:8]
	return
}

// Audit Hooks
func CreatePersonAudit(p Person) (audit PersonAudit) {
	audit.Id = p.Id
	audit.Name = p.Name
	audit.Age = p.Age
	audit.Gender = p.Gender

	return
}

func (p *Person) AfterCreate(tx *gorm.DB) (err error) {
	audit := CreatePersonAudit(*p)
	audit.Operation = "CREATE"
	audit.CreatedBy = "AMAN"
	audit.CreatedAt = time.Now()

	return tx.Create(&audit).Error
}

func (p *Person) AfterUpdate(tx *gorm.DB) (err error) {
	audit := CreatePersonAudit(*p)
	audit.Operation = "UPDATE"
	audit.CreatedBy = "AMAN"
	audit.CreatedAt = time.Now()

	return tx.Create(&audit).Error
}

func (p *Person) AfterDelete(tx *gorm.DB) (err error) {
	audit := CreatePersonAudit(*p)
	audit.Operation = "DELETE"
	audit.CreatedBy = "AMAN"
	audit.CreatedAt = time.Now()

	return tx.Create(&audit).Error
}

// No embedding to decopule Audit and Person
// Also causes issue during save with save loops
type PersonAudit struct {
	Id     string `gorm:"not null"`
	Name   string `gorm:"not null"`
	Age    int    `gorm:"not null"`
	Gender string `gorm:"not null"`

	// Audit Fields
	AuditID   uint      `gorm:"primaryKey"`
	Operation string    `gorm:"not null"`
	CreatedBy string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}
