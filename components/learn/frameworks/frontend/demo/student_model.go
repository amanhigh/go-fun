package main

import "time"

// Student represents a student entity with simplified fields for demo.
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

// StudentListResponse is the paginated student list payload returned by the API.
type StudentListResponse struct {
	Success    bool      `json:"success"`
	Data       []Student `json:"data"`
	Count      int       `json:"count"`
	Offset     int       `json:"offset"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}

// StudentListQuery represents the student list query parameters.
type StudentListQuery struct {
	Offset      int    `form:"offset,default=0"`
	Limit       int    `form:"limit,default=4"`
	SearchQuery string `form:"search"`
	Grade       string `form:"grade"`
}
