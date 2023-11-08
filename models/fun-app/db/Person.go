package db

import "regexp"

var nameRegex = regexp.MustCompile("^[a-zA-Z0-9_-]{1,25}$")

type Person struct {
	Id   int64  `gorm:"primaryKey"`
	Name string `gorm:"not null" binding:"required" validate:"regexp=nameRegex"`
	Age  int    `gorm:"not null" binding:"required" validate:"min=1,max=150"`

	Gender string `gorm:"not null" binding:"required,eq=MALE|eq=FEMALE" enums:"MALE,FEMALE"`

	//TODO: Implement Versioning
	Version int64 `gorm:"not null" json:"-" binding:"-"`
}
