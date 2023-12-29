package handlers

import (
	"net/http"

	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	Manager          manager.PersonManagerInterface `inject:""`
	Tracer           trace.Tracer                   `inject:""`
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
// @Param request body fun.PersonRequest true "Person Request"
// @Success 200 {string} string "Id of created person"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /person [post]
func (self *PersonHandler) CreatePerson(c *gin.Context) {
	/* Captures Create Person Latency */
	timer := prometheus.NewTimer(self.PersonCreateTime)
	defer timer.ObserveDuration()

	ctx, span := self.Tracer.Start(c.Request.Context(), "CreatePerson.Handler")
	defer span.End()

	//Unmarshal the request
	var request fun.PersonRequest
	if err := c.ShouldBind(&request); err == nil {

		self.CreateCounter.WithLabelValues(request.Gender).Inc()

		if person, err := self.Manager.CreatePerson(ctx, request); err == nil {
			c.JSON(http.StatusOK, person)
			span.SetStatus(codes.Ok, "Person Created")
		} else {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			c.JSON(err.Code(), err)
		}
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// GetPerson godoc
//
// @Summary Get a person by ID
// @Description Get a person's details by their ID
// @Tags Person
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Success 200 {object} fun.Person
// @Failure 500 {string} string "Internal Server Error"
// @Router /person/{id} [get]
func (self *PersonHandler) GetPerson(c *gin.Context) {
	var path fun.PersonPath

	ctx, span := self.Tracer.Start(c.Request.Context(), "GetPerson.Handler", trace.WithAttributes(attribute.String("id", path.Id)))
	defer span.End()

	if err := c.ShouldBindUri(&path); err == nil {
		if person, err := self.Manager.GetPerson(ctx, path.Id); err == nil {
			c.JSON(http.StatusOK, person)
		} else {
			c.JSON(err.Code(), err)
		}
	}
}

// ListPersons godoc
//
// @Summary List Person and Search
// @Description List Person and Optionally Search
// @Tags Person
// @Accept json
// @Produce json
// @Param gender query string false "Filter persons by gender"
// @Param age query int false "Filter persons by age"
// @Success 200 {object} fun.PersonList
// @Failure 500 {string} string "Internal Server Error"
// @Router /person [get]
func (self *PersonHandler) ListPersons(c *gin.Context) {
	//FIXME: #C Implement Sorting asc,dsc
	var personQuery fun.PersonQuery

	ctx, span := self.Tracer.Start(c.Request.Context(), "ListPersons.Handler")
	defer span.End()

	if err := c.ShouldBindQuery(&personQuery); err == nil {
		if personList, err := self.Manager.ListPersons(ctx, personQuery); err == nil {
			self.PersonCounter.Add(float64(len(personList.Records)))
			c.JSON(http.StatusOK, personList)
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UpdatePerson godoc
//
// @Summary Update a person
// @Description Update a person's details
// @Tags Person
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Param request body fun.PersonRequest true "Person Request"
// @Success 200 {string} string "UPDATED"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /person/{id} [put]
func (self *PersonHandler) UpdatePerson(c *gin.Context) {
	//https://stackoverflow.com/a/37544666/173136

	ctx, span := self.Tracer.Start(c.Request.Context(), "UpdatePerson.Handler")
	defer span.End()

	//Unmarshal the request
	var request fun.PersonRequest
	if err := c.ShouldBind(&request); err == nil {
		if err := self.Manager.UpdatePerson(ctx, c.Param("id"), request); err == nil {
			//https://stackoverflow.com/a/827045/173136
			c.JSON(http.StatusOK, "UPDATED")
		} else {
			c.JSON(err.Code(), err)
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
	ctx, span := self.Tracer.Start(c.Request.Context(), "DeletePersons.Handler")
	defer span.End()

	if err := self.Manager.DeletePerson(ctx, c.Param("id")); err == nil {
		c.JSON(http.StatusOK, "DELETED")
	} else {
		c.JSON(err.Code(), err)
	}
}
