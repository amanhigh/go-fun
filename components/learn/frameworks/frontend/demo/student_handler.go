package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type studentListResponse struct {
	Success    bool      `json:"success"`
	Data       []Student `json:"data"`
	Count      int       `json:"count"`
	Offset     int       `json:"offset"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}

type studentListQuery struct {
	offset      int
	limit       int
	searchQuery string
	grade       string
}

func validateStudentPayload(student Student) error {
	if student.Age <= 0 {
		return fmt.Errorf("age must be greater than 0")
	}

	return nil
}

// StudentHandler handles student-related HTTP requests
type StudentHandler struct {
	studentService StudentService
}

// NewStudentHandler creates a new student handler
func NewStudentHandler() *StudentHandler {
	return &StudentHandler{
		studentService: NewInMemoryStudentService(),
	}
}

// RegisterRoutes registers student API routes with the Gin router
func (h *StudentHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/students")
	{
		api.GET("", h.getAllStudents)
		api.GET("/:id", h.getStudentByID)
		api.POST("", h.createStudent)
		api.PUT("/:id", h.updateStudent)
		api.DELETE("/:id", h.deleteStudent)
	}
}

// getAllStudents returns a paginated student list as JSON.
func (h *StudentHandler) getAllStudents(c *gin.Context) {
	query := readStudentListQuery(c)
	students, count := h.studentService.GetAllStudents(query.offset, query.limit, query.searchQuery, query.grade)
	writeStudentListResponse(c, students, count, query.offset, query.limit)
}

func readStudentListQuery(c *gin.Context) studentListQuery {
	query := studentListQuery{
		offset: 0,
		limit:  4,
	}

	if value := c.Query("offset"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed >= 0 {
			query.offset = parsed
		}
	}
	if value := c.Query("limit"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			query.limit = parsed
		}
	}
	query.searchQuery = c.Query("search")
	query.grade = c.Query("grade")

	return query
}

func writeStudentListResponse(c *gin.Context, students []Student, count, offset, limit int) {
	if offset >= count && count > 0 {
		offset = max(0, count-limit)
	}
	totalPages := count / limit
	if count%limit != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, studentListResponse{
		Success:    true,
		Data:       students,
		Count:      count,
		Offset:     offset,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// getStudentByID returns a specific student by ID
func (h *StudentHandler) getStudentByID(c *gin.Context) {
	id := c.Param("id")
	student := h.studentService.GetStudentByID(id)

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Student not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    student,
	})
}

// createStudent creates a new student
func (h *StudentHandler) createStudent(c *gin.Context) {
	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := validateStudentPayload(student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	createdStudent := h.studentService.CreateStudent(student)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    createdStudent,
	})
}

// updateStudent updates an existing student
func (h *StudentHandler) updateStudent(c *gin.Context) {
	id := c.Param("id")
	var updatedStudent Student

	if err := c.ShouldBindJSON(&updatedStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := validateStudentPayload(updatedStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	student := h.studentService.UpdateStudent(id, updatedStudent)
	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Student not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    student,
	})
}

// deleteStudent deletes a student by ID
func (h *StudentHandler) deleteStudent(c *gin.Context) {
	id := c.Param("id")

	if !h.studentService.DeleteStudent(id) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Student not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Student deleted successfully",
	})
}
