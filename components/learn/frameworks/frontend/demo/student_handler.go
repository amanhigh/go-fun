package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		api.GET("", h.listStudents)
		api.GET("/:id", h.getStudentByID)
		api.POST("", h.createStudent)
		api.PUT("/:id", h.updateStudent)
		api.DELETE("/:id", h.deleteStudent)
	}
}

// listStudents returns a paginated student list as JSON.
func (h *StudentHandler) listStudents(c *gin.Context) {
	var query StudentListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	response := h.studentService.ListStudents(query.Offset, query.Limit, query.SearchQuery, query.Grade)
	c.JSON(http.StatusOK, response)
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
