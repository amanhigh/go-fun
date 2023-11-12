package server

import (
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun-app/db"
)

type PersonRequest struct {
	db.Person
}

type PersonPath struct {
	Id string `uri:"id" binding:"required"`
}

type PersonQuery struct {
	common.Pagination
	Name   string `binding:"omitempty,min=1,max=25,name=person"`
	Gender string `binding:"omitempty,eq=MALE|eq=FEMALE" enums:"MALE,FEMALE"`
}

type PersonList struct {
	Records []db.Person

	common.PaginatedResponse
}
