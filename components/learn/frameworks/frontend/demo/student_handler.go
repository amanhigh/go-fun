package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

// getAllStudents returns all students as JSON
func (h *StudentHandler) getAllStudents(c *gin.Context) {
	students := h.studentService.GetAllStudents()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    students,
		"count":   len(students),
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
	// FIXME: Add Server Side Validation for one field.
	if err := c.ShouldBindJSON(&student); err != nil {
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
