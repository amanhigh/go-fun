package main

import (
	"fmt"
	"sync"
	"time"
)

// Student represents a student entity with simplified fields for demo
type Student struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Grade     string    `json:"grade"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StudentService interface for student operations
type StudentService interface {
	GetAllStudents() []Student
	GetStudentByID(id string) *Student
	CreateStudent(student Student) *Student
	UpdateStudent(id string, student Student) *Student
	DeleteStudent(id string) bool
}

// InMemoryStudentService implements StudentService with thread-safe in-memory storage
type InMemoryStudentService struct {
	mu       sync.RWMutex
	students map[string]Student
	nextID   int
}

// NewInMemoryStudentService creates a new in-memory student service with sample data
func NewInMemoryStudentService() *InMemoryStudentService {
	service := &InMemoryStudentService{
		students: make(map[string]Student),
		nextID:   1,
	}

	// Add sample students
	sampleStudents := []Student{
		{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@school.edu",
			Age:       20,
			Grade:     "Sophomore",
		},
		{
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@school.edu",
			Age:       21,
			Grade:     "Junior",
		},
		{
			FirstName: "Mike",
			LastName:  "Johnson",
			Email:     "mike.johnson@school.edu",
			Age:       19,
			Grade:     "Freshman",
		},
		{
			FirstName: "Sarah",
			LastName:  "Williams",
			Email:     "sarah.williams@school.edu",
			Age:       22,
			Grade:     "Senior",
		},
		{
			FirstName: "David",
			LastName:  "Brown",
			Email:     "david.brown@school.edu",
			Age:       20,
			Grade:     "Sophomore",
		},
	}

	for _, student := range sampleStudents {
		service.CreateStudent(student)
	}

	return service
}

func (s *InMemoryStudentService) GetAllStudents() []Student {
	s.mu.RLock()
	defer s.mu.RUnlock()

	students := make([]Student, 0, len(s.students))
	for _, student := range s.students {
		students = append(students, student)
	}
	return students
}

func (s *InMemoryStudentService) GetStudentByID(id string) *Student {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if student, exists := s.students[id]; exists {
		return &student
	}
	return nil
}

func (s *InMemoryStudentService) CreateStudent(student Student) *Student {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("%d", s.nextID)
	student.ID = id
	student.CreatedAt = time.Now()
	student.UpdatedAt = time.Now()
	s.students[id] = student
	s.nextID++
	return &student
}

func (s *InMemoryStudentService) UpdateStudent(id string, updatedStudent Student) *Student {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.students[id]; !exists {
		return nil
	}

	updatedStudent.ID = id
	updatedStudent.UpdatedAt = time.Now()
	s.students[id] = updatedStudent
	return &updatedStudent
}

func (s *InMemoryStudentService) DeleteStudent(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.students[id]; !exists {
		return false
	}
	delete(s.students, id)
	return true
}
