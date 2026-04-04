package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
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
	GetAllStudents(offset, limit int, searchQuery, grade string) ([]Student, int)
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
	for _, s := range sampleStudents() {
		service.CreateStudent(s)
	}
	return service
}

// sampleStudents returns seed data for the in-memory store.
func sampleStudents() []Student {
	return []Student{
		{FirstName: "John", LastName: "Doe", Email: "john.doe@school.edu", Age: 20, Grade: "Sophomore"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane.smith@school.edu", Age: 21, Grade: "Junior"},
		{FirstName: "Mike", LastName: "Johnson", Email: "mike.johnson@school.edu", Age: 19, Grade: "Freshman"},
		{FirstName: "Sarah", LastName: "Williams", Email: "sarah.williams@school.edu", Age: 22, Grade: "Senior"},
		{FirstName: "David", LastName: "Brown", Email: "david.brown@school.edu", Age: 20, Grade: "Sophomore"},
		{FirstName: "Emma", LastName: "Davis", Email: "emma.davis@school.edu", Age: 18, Grade: "Freshman"},
		{FirstName: "Liam", LastName: "Miller", Email: "liam.miller@school.edu", Age: 23, Grade: "Senior"},
		{FirstName: "Olivia", LastName: "Wilson", Email: "olivia.wilson@school.edu", Age: 20, Grade: "Junior"},
		{FirstName: "Noah", LastName: "Moore", Email: "noah.moore@school.edu", Age: 19, Grade: "Sophomore"},
		{FirstName: "Ava", LastName: "Taylor", Email: "ava.taylor@school.edu", Age: 21, Grade: "Senior"},
		{FirstName: "Ethan", LastName: "Anderson", Email: "ethan.anderson@school.edu", Age: 18, Grade: "Freshman"},
		{FirstName: "Sophia", LastName: "Thomas", Email: "sophia.thomas@school.edu", Age: 22, Grade: "Junior"},
		{FirstName: "Lucas", LastName: "Jackson", Email: "lucas.jackson@school.edu", Age: 20, Grade: "Sophomore"},
		{FirstName: "Mia", LastName: "White", Email: "mia.white@school.edu", Age: 19, Grade: "Freshman"},
		{FirstName: "Mason", LastName: "Harris", Email: "mason.harris@school.edu", Age: 23, Grade: "Senior"},
		{FirstName: "Isabella", LastName: "Martin", Email: "isabella.martin@school.edu", Age: 20, Grade: "Junior"},
		{FirstName: "Logan", LastName: "Thompson", Email: "logan.thompson@school.edu", Age: 18, Grade: "Freshman"},
		{FirstName: "Amelia", LastName: "Garcia", Email: "amelia.garcia@school.edu", Age: 21, Grade: "Senior"},
		{FirstName: "Elijah", LastName: "Martinez", Email: "elijah.martinez@school.edu", Age: 19, Grade: "Sophomore"},
		{FirstName: "Harper", LastName: "Robinson", Email: "harper.robinson@school.edu", Age: 22, Grade: "Junior"},
	}
}

func (s *InMemoryStudentService) GetAllStudents(offset, limit int, searchQuery, grade string) ([]Student, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	students := s.filteredStudents(searchQuery, grade)
	total := len(students)
	if limit <= 0 {
		limit = total
	}
	if limit <= 0 {
		return students, total
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= total {
		offset = max(0, total-limit)
	}
	return sliceOffset(students, offset, limit), total
}

func (s *InMemoryStudentService) filteredStudents(searchQuery, grade string) []Student {
	students := s.sortedStudents()
	query := strings.ToLower(strings.TrimSpace(searchQuery))
	grade = strings.TrimSpace(grade)
	if query == "" && grade == "" {
		return students
	}

	filtered := make([]Student, 0, len(students))
	for _, student := range students {
		if !matchesStudentFilters(student, query, grade) {
			continue
		}
		filtered = append(filtered, student)
	}
	return filtered
}

func matchesStudentFilters(student Student, query, grade string) bool {
	if grade != "" && student.Grade != grade {
		return false
	}
	if query == "" {
		return true
	}
	name := strings.ToLower(student.FirstName + " " + student.LastName)
	return strings.Contains(name, query)
}

func (s *InMemoryStudentService) sortedStudents() []Student {
	students := make([]Student, 0, len(s.students))
	for _, student := range s.students {
		students = append(students, student)
	}
	sort.Slice(students, func(i, j int) bool {
		left, leftErr := strconv.Atoi(students[i].ID)
		right, rightErr := strconv.Atoi(students[j].ID)
		if leftErr != nil || rightErr != nil {
			return students[i].ID < students[j].ID
		}
		return left < right
	})
	return students
}

func sliceOffset(students []Student, offset, limit int) []Student {
	if offset >= len(students) {
		return []Student{}
	}
	end := min(offset+limit, len(students))
	return students[offset:end]
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
