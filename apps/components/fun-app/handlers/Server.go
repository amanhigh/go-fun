package handlers

import (
	"fmt"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type FunServer struct {
	GinEngine *gin.Engine `inject:""`
	Port      int         `inject:"port"`

	/* Handlers */
	PersonHandlers PersonHandler `inject:"inline"`
}

func (self *FunServer) initRoutes() {
	//Routes
	personGroup := self.GinEngine.Group("/person")
	personGroup.GET("/all", self.PersonHandlers.GetAllPerson)
	personGroup.POST("", self.PersonHandlers.CreatePerson)
	personGroup.DELETE(":id", self.PersonHandlers.DeletePersons)

	//Pprof (Use: http://localhost:8080/debug/pprof/)
	pprof.Register(self.GinEngine)
}

func (self *FunServer) Start() (err error) {
	self.initRoutes()
	return self.GinEngine.Run(fmt.Sprintf(":%v", self.Port))
}
