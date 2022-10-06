package handlers

import (
	"net/http"

	manager2 "github.com/amanhigh/go-fun/components/fun-app/manager"
	server2 "github.com/amanhigh/go-fun/models/fun-app/server"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	Manager          manager2.PersonManagerInterface `inject:""`
	CreateCounter    *prometheus.CounterVec          `inject:"m_create_person"`
	PersonCounter    prometheus.Gauge                `inject:"m_person_count"`
	PersonCreateTime prometheus.Histogram            `inject:"m_person_create_time"`
}

func (self *PersonHandler) CreatePerson(c *gin.Context) {
	/* Captures Create Person Latency */
	timer := prometheus.NewTimer(self.PersonCreateTime)
	defer timer.ObserveDuration()

	var request server2.PersonRequest
	if err := c.Bind(&request); err == nil {
		self.CreateCounter.WithLabelValues(request.Gender).Inc()

		if err := self.Manager.CreatePerson(c, request.Person); err == nil {
			c.JSON(http.StatusOK, request)
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

func (self *PersonHandler) GetPerson(c *gin.Context) {
	if person, err := self.Manager.GetPerson(c, c.Param("id")); err == nil {
		c.JSON(http.StatusOK, person)
	} else {
		c.JSON(http.StatusInternalServerError, err.Error())
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
