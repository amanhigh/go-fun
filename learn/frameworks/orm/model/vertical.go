package model

//go:generate jsongen -type=Vertical -package model

import "github.com/jinzhu/gorm"

type Vertical struct {
	gorm.Model
	Name     string `gorm:"unique;default:'Shirts'"`
	MyColumn string
}
