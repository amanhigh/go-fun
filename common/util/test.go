package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

const DIMENSION_COUNT = 2 // Number of values needed for matrix dimensions (n, m)

func ReadCountInts(scanner *bufio.Scanner) (n int, ints []int) {
	n = ReadInt(scanner)
	ints = ReadInts(scanner, n)
	return
}

func ReadMatrix(scanner *bufio.Scanner, n, m int) (matrix [][]int) {
	matrix = make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = ReadInts(scanner, m)
	}
	return
}

func ReadMatrixWithDimensions(scanner *bufio.Scanner) (matrix [][]int, n, m int) {
	ints := ReadInts(scanner, DIMENSION_COUNT)
	n = ints[0]
	m = ints[1]
	matrix = ReadMatrix(scanner, n, m)
	return
}

// CreateTestRequest creates an HTTP request for testing with JSON headers
func CreateTestRequest(method, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var jsonData []byte
	var err error

	if body != nil {
		jsonData, err = json.Marshal(body)
		if err != nil {
			log.Fatalf("Failed to marshal Test Request JSON: %v", err)
		}
	} else {
		jsonData = []byte{}
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	return req, w
}

// AssertJSONAndStatus combines JSON unmarshaling and status checking
func AssertJSONAndStatus(w *httptest.ResponseRecorder, expectedStatus int, v any) {
	if w.Code != expectedStatus {
		log.Fatalf("Expected status %d, got %d", expectedStatus, w.Code)
	}
	if err := json.Unmarshal(w.Body.Bytes(), v); err != nil {
		log.Fatalf("Failed to unmarshal Test JSON: %v", err)
	}
}
