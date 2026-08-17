package handlers

import (
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
)

type StudentHandler interface {
	CreateStudent(c *gin.Context)
	GetStudent(c *gin.Context)
	ListStudents(c *gin.Context)
	ListStudentAudit(c *gin.Context)
	UpdateStudent(c *gin.Context)
	DeleteStudents(c *gin.Context)
}

type StudentHandlerImpl struct {
	Manager          manager.StudentManagerInterface `container:"type"`
	Tracer           trace.Tracer                   `container:"type"`
	CreateCounter    metric.Int64Counter            `container:"name"`
	StudentCounter    metric.Int64UpDownCounter      `container:"name"`
	StudentCreateTime metric.Float64Histogram        `container:"name"`
}

// CreateStudent godoc
//
// @Summary Create a new student
// @Description Create a new student with the provided data
// @Tags Student
// @Accept json
// @Produce json
// @Param request body fun.StudentRequest true "Student Request"
// @Success 201 {string} string "Id of created student"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /student [post]
func (ph *StudentHandlerImpl) CreateStudent(c *gin.Context) {
	/* Captures Create Student Latency */
	startTime := time.Now()
	defer func() {
		ph.StudentCreateTime.Record(c.Request.Context(), time.Since(startTime).Seconds())
	}()

	ctx, span := ph.Tracer.Start(c.Request.Context(), "CreateStudent.Handler")
	defer span.End()

	// Unmarshal the request
	var request fun.StudentRequest
	if err := c.ShouldBind(&request); err == nil {
		ph.CreateCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("gender", request.Gender)))

		// FIXME: Move Student Flows and API's to Envelope similar to kohan API's
		if student, err := ph.Manager.CreateStudent(ctx, request); err == nil {
			c.JSON(http.StatusCreated, student)
			span.SetStatus(codes.Ok, "Student Created")
		} else {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			c.JSON(err.Code(), err)
		}
	} else {
		httpErr := util.ProcessValidationError(err)
		c.JSON(httpErr.Code(), httpErr)
	}
}

// GetStudent godoc
//
// @Summary Get a student by ID
// @Description Get a student's details by their ID
// @Tags Student
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} fun.Student
// @Failure 500 {string} string "Internal Server Error"
// @Router /student/{id} [get]
func (ph *StudentHandlerImpl) GetStudent(c *gin.Context) {
	var path fun.StudentPath

	ctx, span := ph.Tracer.Start(c.Request.Context(), "GetStudent.Handler", trace.WithAttributes(attribute.String("id", path.Id)))
	defer span.End()

	util.UnwrapRequest(c, nil, &path, nil, func(c *gin.Context) {
		if student, err := ph.Manager.GetStudent(ctx, path.Id); err == nil {
			c.JSON(http.StatusOK, student)
		} else {
			httpErr := util.ProcessValidationError(err)
			span.SetStatus(codes.Error, httpErr.Error())
			span.RecordError(httpErr)
			c.JSON(httpErr.Code(), httpErr)
		}
	})
}

// ListStudents godoc
//
// @Summary List Student and Search
// @Description List Student and Optionally Search
// @Tags Student
// @Accept json
// @Produce json
// @Param name query string false "Filter students by name"
// @Param gender query string false "Filter students by gender"
// @Param age query int false "Filter students by age"
// @Param order query string false "Sort order" Enums(asc, desc)
// @Param sort_by query string false "Sort by" Enums(name, gender, age)
// @Success 200 {object} fun.StudentList
// @Failure 500 {string} string "Internal Server Error"
// @Router /student [get]
func (ph *StudentHandlerImpl) ListStudents(c *gin.Context) {
	var studentQuery fun.StudentQuery

	ctx, span := ph.Tracer.Start(c.Request.Context(), "ListStudents.Handler")
	defer span.End()

	if err := c.ShouldBindQuery(&studentQuery); err == nil {
		if studentList, err := ph.Manager.ListStudents(ctx, studentQuery); err == nil {
			count := int64(len(studentList.Records))
			ph.StudentCounter.Add(ctx, count)
			c.JSON(http.StatusOK, studentList)
		} else {
			zerolog.Ctx(ctx).Error().Err(err).Msg("ListStudents: Server Error")
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		httpErr := util.ProcessValidationError(err)
		zerolog.Ctx(ctx).Error().Err(httpErr).Msg("ListStudents: Bad Request")
		c.JSON(httpErr.Code(), httpErr)
	}
}

// ListStudentAudit godoc
//
// @Summary List Student Audit
// @Description List Student Audit by ID
// @Tags Student
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} []fun.StudentAudit
// @Failure 500 {string} string "Internal Server Error"
// @Router /student/{id}/audit [get]
func (ph *StudentHandlerImpl) ListStudentAudit(c *gin.Context) {
	var path fun.StudentPath

	ctx, span := ph.Tracer.Start(c.Request.Context(), "ListStudentAudit.Handler", trace.WithAttributes(attribute.String("id", path.Id)))
	defer span.End()

	if err := c.ShouldBindUri(&path); err == nil {
		if auditList, err := ph.Manager.ListStudentAudit(ctx, path.Id); err == nil {
			c.JSON(http.StatusOK, auditList)
		} else {
			httpErr := util.ProcessValidationError(err)
			zerolog.Ctx(ctx).Err(httpErr).Int("status", httpErr.Code()).Msg("ListStudentAudit: Server Error")
			c.JSON(httpErr.Code(), httpErr)
		}
	}
}

// UpdateStudent godoc
//
// @Summary Update a student
// @Description Update a student's details
// @Tags Student
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body fun.StudentRequest true "Student Request"
// @Success 200 {string} string "UPDATED"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /student/{id} [put]
func (ph *StudentHandlerImpl) UpdateStudent(c *gin.Context) {
	//https://stackoverflow.com/a/37544666/173136

	ctx, span := ph.Tracer.Start(c.Request.Context(), "UpdateStudent.Handler")
	defer span.End()

	// Unmarshal the request
	var request fun.StudentRequest
	if err := c.ShouldBind(&request); err == nil {
		if err := ph.Manager.UpdateStudent(ctx, c.Param("id"), request); err == nil {
			//https://stackoverflow.com/a/827045/173136
			c.JSON(http.StatusOK, "UPDATED")
		} else {
			zerolog.Ctx(ctx).Err(err).Int("status", err.Code()).Msg("UpdateStudent: Server Error")
			c.JSON(err.Code(), err)
		}
	} else {
		httpErr := util.ProcessValidationError(err)
		zerolog.Ctx(ctx).Err(httpErr).Int("status", http.StatusBadRequest).Msg("UpdateStudent: Bad Request")
		c.JSON(httpErr.Code(), httpErr)
	}
}

// DeleteStudents godoc
//
// @Summary Delete students by ID
// @Description Delete students by their ID
// @Tags Student
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {string} string "DELETED"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /student/{id} [delete]
func (ph *StudentHandlerImpl) DeleteStudents(c *gin.Context) {
	ctx, span := ph.Tracer.Start(c.Request.Context(), "DeleteStudents.Handler")
	defer span.End()

	if err := ph.Manager.DeleteStudent(ctx, c.Param("id")); err == nil {
		c.JSON(http.StatusNoContent, "DELETED")
	} else {
		httpErr := util.ProcessValidationError(err)
		zerolog.Ctx(ctx).Err(httpErr).Int("status", httpErr.Code()).Msg("DeleteStudents: Server Error")
		c.JSON(httpErr.Code(), httpErr)
	}
}
