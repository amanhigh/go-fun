package server

import (
	db2 "github.com/amanhigh/go-fun/models/fun-app/db"
)

type PersonRequest struct {
	db2.Person
}

type PersonPath struct {
	Id string `uri:"id" binding:"required"`
}
