package fun

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const CreatedByAman = "AMAN"

type StudentRequest struct {
	// Validations - https://gin-gonic.com/docs/examples/binding-and-validation/
	Name   string `json:"name" gorm:"not null" binding:"required,min=1,max=25,name=student"`
	Age    int    `json:"age" gorm:"not null" binding:"required,min=1,max=150"`
	Gender string `json:"gender" gorm:"not null" binding:"required,eq=MALE|eq=FEMALE" enums:"MALE,FEMALE"`
}

type StudentPath struct {
	Id string `uri:"id" binding:"required"`
}

type StudentQuery struct {
	common.Pagination
	common.Sort
	Name   string `form:"name" binding:"omitempty,min=1,max=25,name=student"`
	Gender string `form:"gender" binding:"omitempty,eq=MALE|eq=FEMALE"`
	SortBy string `form:"sort_by" binding:"omitempty,eq=name|eq=age|eq=gender"`
}

type StudentList struct {
	Records  []Student                `json:"records"`
	Metadata common.PaginatedResponse `json:"metadata"`
}

type Student struct {
	StudentRequest
	Id string `gorm:"primaryKey" json:"id"`
}

func (p *Student) BeforeCreate(_ *gorm.DB) (err error) {
	p.Id = uuid.NewString()[:8]
	return
}

// Audit Hooks
func CreateStudentAudit(p Student) (audit StudentAudit) {
	audit.Id = p.Id
	audit.Name = p.Name
	audit.Age = p.Age
	audit.Gender = p.Gender

	return
}

func (p *Student) AfterCreate(tx *gorm.DB) (err error) {
	audit := CreateStudentAudit(*p)
	audit.Operation = "CREATE"
	audit.CreatedBy = CreatedByAman
	audit.CreatedAt = time.Now()

	return tx.Create(&audit).Error
}

func (p *Student) AfterUpdate(tx *gorm.DB) (err error) {
	audit := CreateStudentAudit(*p)
	audit.Operation = "UPDATE"
	audit.CreatedBy = CreatedByAman
	audit.CreatedAt = time.Now()

	return tx.Create(&audit).Error
}

func (p *Student) AfterDelete(tx *gorm.DB) (err error) {
	audit := CreateStudentAudit(*p)
	audit.Operation = "DELETE"
	audit.CreatedBy = CreatedByAman
	audit.CreatedAt = time.Now()

	return tx.Create(&audit).Error
}

// No embedding to decopule Audit and Student
// Also causes issue during save with save loops
type StudentAudit struct {
	Id     string `gorm:"not null"`
	Name   string `gorm:"not null"`
	Age    int    `gorm:"not null"`
	Gender string `gorm:"not null"`

	// Audit Fields
	AuditID   uint   `gorm:"primaryKey"`
	Operation string `gorm:"not null"`
	// HACK: Use Base Dao of Gorm for common Fields ?
	CreatedBy string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}
