package frameworks

import (
	"encoding/json"

	"gorm.io/gorm"
)

// Field Tags - https://gorm.io/docs/models.html#Fields-Tags
type Product struct {
	gorm.Model
	Code       string `gorm:"size 5,unique"`
	Price      uint   `gorm:"not null"`
	Version    int
	IgnoreMe   string `gorm:"-"` // Ignore this field
	Vertical   Vertical
	VerticalID uint      // Must be vertical_id in DB or won't work automatically.
	Features   []Feature `gorm:"many2many:product_features;"`
}

type AuditLog struct {
	gorm.Model
	Operation string
	Log       string
}

type Feature struct {
	gorm.Model
	Name    string
	Version int
}

type Vertical struct {
	gorm.Model
	Name     string `gorm:"unique;default:'Shirts'"`
	MyColumn string
}

// Default Name would be products
func (p *Product) TableName() string {
	return "MeraProduct"
}

// begin transaction
// -> BeforeSave
// -> BeforeCreate/Update
// save before associations
// update timestamp `CreatedAt`, `UpdatedAt`
// save self
// reload fields that have default value and its value is blank
// save after associations
// -> AfterCreate
// -> AfterSave/Update
// commit or rollback transaction

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	// Log Product
	marshal, _ := json.Marshal(p)
	p.Version++
	tx.Create(&AuditLog{Operation: "Create", Log: string(marshal)})
	return
}

func (p *Product) BeforeUpdate(tx *gorm.DB) (err error) {
	// Log Product
	marshal, _ := json.Marshal(p)
	p.Version++
	tx.Create(&AuditLog{Operation: "Update", Log: string(marshal)})
	return
}

func (f *Feature) BeforeCreate(tx *gorm.DB) (err error) {
	// Log Feature
	marshal, _ := json.Marshal(f)
	f.Version++
	tx.Create(&AuditLog{Operation: "Create", Log: string(marshal)})
	return
}

// Use Value instead of pointer for delete as no version update is required
func (f Feature) BeforeDelete(tx *gorm.DB) (err error) {
	// Log Feature
	marshal, _ := json.Marshal(f)
	tx.Create(&AuditLog{Operation: "Delete", Log: string(marshal)})
	return
}

func (p *Product) AfterFind(_ *gorm.DB) (err error) {
	p.IgnoreMe = "Ignore" + p.Code
	return nil
}
