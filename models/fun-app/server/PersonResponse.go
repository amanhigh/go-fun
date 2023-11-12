package server

import (
	"github.com/amanhigh/go-fun/models/common"
	db "github.com/amanhigh/go-fun/models/fun-app/db"
)

type PersonQuery struct {
	common.Pagination
	Name   string
	Gender string
}

type PersonList struct {
	Records []db.Person

	common.PaginatedResponse
}
