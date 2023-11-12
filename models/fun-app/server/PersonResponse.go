package server

import (
	"github.com/amanhigh/go-fun/models/common"
	db "github.com/amanhigh/go-fun/models/fun-app/db"
)

type PersonList struct {
	Records []db.Person

	common.PaginatedResponse
}
