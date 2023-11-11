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

// CreatePerson godoc
//
// @Summary Create a new person
// @Description Create a new person with the provided data
// @Tags Person
// @Accept json
// @Produce json
// @Param request body server2.PersonRequest true "Person Request"
// @Success 200 {string} string "Id of created person"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /person [post]
func (self *PersonHandler) CreatePerson(c *gin.Context) {
	/* Captures Create Person Latency */
	timer := prometheus.NewTimer(self.PersonCreateTime)
	defer timer.ObserveDuration()

	var request server2.PersonRequest
	if err := c.ShouldBind(&request); err == nil {

		self.CreateCounter.WithLabelValues(request.Gender).Inc()

		if id, err := self.Manager.CreatePerson(c, request.Person); err == nil {
			c.JSON(http.StatusOK, id)
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// GetPerson godoc
//
//	@Summary		Get a person by ID
//	@Description	Get a person's details by their ID
//	@Tags			Person
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Person ID"
//	@Success		200	{object}	db.Person
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/person/{id} [get]
func (self *PersonHandler) GetPerson(c *gin.Context) {
	if person, err := self.Manager.GetPerson(c, c.Param("id")); err == nil {
		c.JSON(http.StatusOK, person)
	} else {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}

// GetAllPerson godoc
//
// @Summary Get all persons
// @Description Get all persons' details
// @Tags Person
// @Accept json
// @Produce json
// @Success 200 {array} db.Person
// @Failure 500 {string} string "Internal Server Error"
// @Router /person/all [get]
func (self *PersonHandler) GetAllPerson(c *gin.Context) {
	if persons, err := self.Manager.GetAllPersons(c); err == nil {
		self.PersonCounter.Add(float64(len(persons)))
		c.JSON(http.StatusOK, persons)
	} else {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}

// DeletePersons godoc
//
// @Summary Delete persons by ID
// @Description Delete persons by their ID
// @Tags Person
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Success 200 {string} string "DELETED"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /person/{id} [delete]
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
