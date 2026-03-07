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

// EnrollmentHandler exposes REST endpoints for enrollment resources.
type EnrollmentHandler interface {
	CreateEnrollment(c *gin.Context)
	GetEnrollment(c *gin.Context)
}

type EnrollmentHandlerImpl struct {
	Manager manager.EnrollmentManagerInterface `container:"type"`
	Tracer  trace.Tracer                       `container:"type"`
}

// NewEnrollmentHandler constructs handler with explicit dependencies and returns interface.
func NewEnrollmentHandler(manager manager.EnrollmentManagerInterface, tracer trace.Tracer) EnrollmentHandler {
	h := &EnrollmentHandlerImpl{Manager: manager, Tracer: tracer}
	return h
}

// CreateEnrollment orchestrates enrollment using an existing person record.
func (eh *EnrollmentHandlerImpl) CreateEnrollment(c *gin.Context) {
	ctx, span := eh.Tracer.Start(c.Request.Context(), "CreateEnrollment.Handler")
	defer span.End()

	var request fun.EnrollmentRequest
	if bindErr := c.ShouldBindJSON(&request); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		span.RecordError(httpErr)
		span.SetStatus(codes.Error, httpErr.Error())
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	enrollment, err := eh.Manager.EnrollPerson(ctx, request.PersonID, request.Grade)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(err.Code(), err)
		return
	}

	c.Header("Location", fmt.Sprintf("/v1/enrollments/%s", enrollment.PersonID))
	span.SetStatus(codes.Ok, "Enrollment accepted")
	c.JSON(http.StatusAccepted, enrollment)
}

// GetEnrollment fetches enrollment status for a person.
func (eh *EnrollmentHandlerImpl) GetEnrollment(c *gin.Context) {
	ctx, span := eh.Tracer.Start(c.Request.Context(), "GetEnrollment.Handler")
	defer span.End()

	var path fun.EnrollmentPath
	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		span.RecordError(httpErr)
		span.SetStatus(codes.Error, httpErr.Error())
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	enrollment, err := eh.Manager.GetEnrollment(ctx, path.PersonID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(err.Code(), err)
		return
	}

	span.SetStatus(codes.Ok, "Enrollment retrieved")
	c.JSON(http.StatusOK, enrollment)
}
