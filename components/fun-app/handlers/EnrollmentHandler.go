package handlers

import (
	"fmt"
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type EnrollmentHandler struct {
	Manager manager.EnrollmentManagerInterface `container:"type"`
	Tracer  trace.Tracer                       `container:"type"`
}

func NewEnrollmentHandler() *EnrollmentHandler {
	return &EnrollmentHandler{}
}

// CreateEnrollment orchestrates enrollment using an existing person record.
func (eh *EnrollmentHandler) CreateEnrollment(c *gin.Context) {
	ctx, span := eh.Tracer.Start(c.Request.Context(), "CreateEnrollment.Handler")
	defer span.End()

	var request fun.EnrollmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpErr := util.ProcessValidationError(err)
		span.RecordError(httpErr)
		span.SetStatus(codes.Error, httpErr.Error())
		c.JSON(http.StatusBadRequest, httpErr)
		return
	}

	person, err := eh.Manager.EnrollPerson(ctx, request.PersonID, request.Grade)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(err.Code(), err)
		return
	}

	response := fun.EnrollmentResponse{
		PersonID: person.Id,
		Grade:    request.Grade,
		Status:   "ACTIVE",
		Links: map[string]string{
			"person": fmt.Sprintf("/v1/person/%s", person.Id),
		},
	}

	span.SetStatus(codes.Ok, "Enrollment completed")
	c.JSON(http.StatusCreated, response)
}
