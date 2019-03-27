package frameworks

import "time"

type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

type DocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DocumentResponse struct {
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

type SearchResponse struct {
	Time      string             `json:"time"`
	Hits      string             `json:"hits"`
	Documents []DocumentResponse `json:"documents"`
}
