package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
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
// FIXME: Add Modernize Linter (Once upgraded to go 1.26) for interface to any checks for golang and other modernization.
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

// UnenvelopeAndAssertStatus combines envelope unwrapping with status checking
// Use this for API responses that return common.Envelope[T] format
func UnenvelopeAndAssertStatus[T any](w *httptest.ResponseRecorder, expectedStatus int) T {
	if w.Code != expectedStatus {
		log.Fatalf("Expected status %d, got %d", expectedStatus, w.Code)
	}

	var envelope common.Envelope[T]
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		log.Fatalf("Failed to unmarshal Envelope JSON: %v", err)
	}

	return envelope.Data
}

// CreateTestGinRouter creates a new Gin router configured for testing
// Used to eliminate duplication across handler integration tests
func CreateTestGinRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}
