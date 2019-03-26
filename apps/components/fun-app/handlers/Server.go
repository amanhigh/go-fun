package handlers

import (
	"fmt"

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
}

func (self *FunServer) Start() (err error) {
	self.initRoutes()
	return self.GinEngine.Run(fmt.Sprintf(":%v", self.Port))
}
