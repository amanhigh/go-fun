package handlers

import (
	"net/http"

	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun-app/server"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	Manager          manager.PersonManagerInterface `inject:""`
	CreateCounter    *prometheus.CounterVec         `inject:"m_create_person"`
	PersonCounter    prometheus.Gauge               `inject:"m_person_count"`
	PersonCreateTime prometheus.Histogram           `inject:"m_person_create_time"`
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

	//Unmarshal the request
	var request server.PersonRequest
	if err := c.ShouldBind(&request); err == nil {

		self.CreateCounter.WithLabelValues(request.Gender).Inc()

		if id, err := self.Manager.CreatePerson(c, request); err == nil {
			c.JSON(http.StatusOK, id)
		} else {
			c.JSON(err.Code(), err)
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
	var path server.PersonPath

	if err := c.ShouldBindUri(&path); err == nil {
		if person, err := self.Manager.GetPerson(c, path.Id); err == nil {
			c.JSON(http.StatusOK, person)
		} else {
			c.JSON(err.Code(), err)
		}
	}

}

// ListPersons godoc
//
// @Summary Get all persons
// @Description Get all persons' details
// @Tags Person
// @Accept json
// @Produce json
// @Success 200 {array} db.Person
// @Failure 500 {string} string "Internal Server Error"
// @Router /person/all [get]
func (self *PersonHandler) ListPersons(c *gin.Context) {
	var pageParams common.Pagination
	if err := c.ShouldBindQuery(&pageParams); err == nil {
		if personList, err := self.Manager.ListPersons(c, pageParams); err == nil {
			self.PersonCounter.Add(float64(len(personList.Records)))
			c.JSON(http.StatusOK, personList)
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.JSON(http.StatusBadRequest, err)
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
		c.JSON(err.Code(), err)
	}
}
