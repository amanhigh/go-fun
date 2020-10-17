package model

import "gorm.io/gorm"

//go:generate jsongen -type=Vertical -package model

type Vertical struct {
	gorm.Model
	Name     string `gorm:"unique;default:'Shirts'"`
	MyColumn string
}
