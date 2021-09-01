package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"net/http"

	"github.com/amanhigh/go-fun/apps/components/fun-app/manager"
	"github.com/amanhigh/go-fun/apps/models/fun-app/server"
	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	Manager          manager.PersonManagerInterface `inject:""`
	CreateCounter    prometheus.Counter             `inject:"m_create_person"`
	PersonCounter    prometheus.Gauge               `inject:"m_person_count"`
	PersonCreateTime prometheus.Histogram           `inject:"m_person_create_time"`
}

func (self *PersonHandler) CreatePerson(c *gin.Context) {
	/* Captures Create Person Latency */
	timer := prometheus.NewTimer(self.PersonCreateTime)
	defer timer.ObserveDuration()

	self.CreateCounter.Inc()

	var request server.PersonRequest
	if err := c.Bind(&request); err == nil {
		if err := self.Manager.CreatePerson(c, request.Person); err == nil {
			c.JSON(http.StatusOK, request)
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

func (self *PersonHandler) GetAllPerson(c *gin.Context) {
	if persons, err := self.Manager.GetAllPersons(c); err == nil {
		self.PersonCounter.Add(float64(len(persons)))
		c.JSON(http.StatusOK, persons)
	} else {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (self *PersonHandler) DeletePersons(c *gin.Context) {
	if err := self.Manager.DeletePerson(c, c.Param("id")); err == nil {
		c.JSON(http.StatusOK, "DELETED")
	} else {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, err)
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
}
