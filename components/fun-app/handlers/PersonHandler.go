package handlers

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type PersonHandler struct {
	Manager          manager.PersonManagerInterface `container:"type"`
	Tracer           trace.Tracer                   `container:"type"`
	CreateCounter    *prometheus.CounterVec         `container:"name"`
	PersonCounter    prometheus.Gauge               `container:"name"`
	PersonCreateTime prometheus.Histogram           `container:"name"`
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
		err = util.ProcessValidationError(err)
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
			c.JSON(http.StatusCreated, person)
		} else {
			err = util.ProcessValidationError(err)
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
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
// @Param name query string false "Filter persons by name"
// @Param gender query string false "Filter persons by gender"
// @Param age query int false "Filter persons by age"
// @Param order query string false "Sort order" Enums(asc, desc)
// @Param sort_by query string false "Sort by" Enums(name, gender, age)
// @Success 200 {object} fun.PersonList
// @Failure 500 {string} string "Internal Server Error"
// @Router /person [get]
func (self *PersonHandler) ListPersons(c *gin.Context) {
	var personQuery fun.PersonQuery
	personQuery.Order = "asc" //Default Sort Order

	ctx, span := self.Tracer.Start(c.Request.Context(), "ListPersons.Handler")
	defer span.End()

	if err := c.ShouldBindQuery(&personQuery); err == nil {
		if personList, err := self.Manager.ListPersons(ctx, personQuery); err == nil {
			self.PersonCounter.Add(float64(len(personList.Records)))
			c.JSON(http.StatusOK, personList)
		} else {
			log.WithFields(log.Fields{"Error": err}).Error("ListPersons: Server Error")
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		err = util.ProcessValidationError(err)
		log.WithFields(log.Fields{"Error": err}).Error("ListPersons: Bad Request")
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
			log.WithFields(log.Fields{"Error": err}).Error("UpdatePerson: Server Error")
		}
	} else {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		log.WithFields(log.Fields{"Error": err}).Error("UpdatePerson: Bad Request")
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
		c.JSON(http.StatusNoContent, "DELETED")
	} else {
		err = util.ProcessValidationError(err)
		c.JSON(err.Code(), err)
		log.WithFields(log.Fields{"Error": err}).Error("DeletePersons: Server Error")
	}
}
