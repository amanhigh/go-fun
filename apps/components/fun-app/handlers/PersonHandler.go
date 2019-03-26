package handlers

import (
	"net/http"

	"github.com/amanhigh/go-fun/apps/components/fun-app/manager"
	"github.com/amanhigh/go-fun/apps/models/fun-app/server"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type PersonHandler struct {
	Manager manager.PersonManagerInterface `inject:""`
}

func (self *PersonHandler) CreatePerson(c *gin.Context) {
	var request server.PersonRequest
	if err := c.Bind(&request); err == nil {
		if err := self.Manager.CreatePerson(request.Person); err == nil {
			c.JSON(http.StatusOK, request)
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

func (self *PersonHandler) GetAllPerson(c *gin.Context) {
	if persons, err := self.Manager.GetAllPersons(); err == nil {
		c.JSON(http.StatusOK, persons)
	} else {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (self *PersonHandler) DeletePersons(c *gin.Context) {
	if err := self.Manager.DeletePerson(c.Param("id")); err == nil {
		c.JSON(http.StatusOK, "DELETED")
	} else {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, err)
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
}
