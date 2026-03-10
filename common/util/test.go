package util

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

const DIMENSION_COUNT = 2 // Number of values needed for matrix dimensions (n, m)

func ReadCountInts(scanner *bufio.Scanner) (n int, ints []int) {
	n = ReadInt(scanner)
	ints = ReadInts(scanner, n)
	return
}

func ReadMatrix(scanner *bufio.Scanner, n, m int) (matrix [][]int) {
	matrix = make([][]int, n)
	for i := range n {
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
func CreateTestRequest(method, url string, body any) (*http.Request, *httptest.ResponseRecorder) {
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

	req := httptest.NewRequestWithContext(context.Background(), method, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	return req, w
}

// AssertSuccess validates 2xx success responses and unmarshals JSON
// Use this for happy path responses (200, 201, 204, etc.)
// Only accepts 2xx status codes and uses Omega assertions without panic
func AssertSuccess(w *httptest.ResponseRecorder, expectedStatus int, v any) {
	// Validate that expectedStatus is a 2xx code using Omega assertions
	ExpectWithOffset(1, expectedStatus).To(BeNumerically(">=", 200), "AssertSuccess only accepts 2xx status codes, got %d", expectedStatus)
	ExpectWithOffset(1, expectedStatus).To(BeNumerically("<", 300), "AssertSuccess only accepts 2xx status codes, got %d", expectedStatus)

	// Check status code matches using Omega assertions
	ExpectWithOffset(1, w.Code).To(Equal(expectedStatus), "Expected status %d, got %d", expectedStatus, w.Code)

	// Unmarshal JSON response using Omega assertions
	err := json.Unmarshal(w.Body.Bytes(), v)
	ExpectWithOffset(1, err).ToNot(HaveOccurred(), "Failed to unmarshal Test JSON: %v", err)
}

// CreateTestGinRouter creates a new Gin router configured for testing
// Used to eliminate duplication across handler integration tests
func CreateTestGinRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// AssertError validates error responses (primarily 4xx field validation errors)
// Uses Omega assertions instead of panic for cleaner test failures
// All field validation errors return 400 Bad Request, so status code is fixed
// Takes field name and expected content to validate, returns nothing
func AssertError(w *httptest.ResponseRecorder, fieldName, expectedContent string) {
	// Check status code is 400
	Expect(w.Code).To(Equal(http.StatusBadRequest))

	// Parse JSON response
	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	Expect(err).ToNot(HaveOccurred())

	// Check JSend format
	Expect(response["status"]).To(Equal("fail"))
	Expect(response["data"]).ToNot(BeNil())

	// Get data map
	dataMap := response["data"].(map[string]any)

	// Check field exists and validate content
	errorMsg, ok := dataMap[fieldName]
	Expect(ok).To(BeTrue(), "expected field '"+fieldName+"' not found in JSend fail data")

	errorMsgStr, ok := errorMsg.(string)
	Expect(ok).To(BeTrue(), "field error message is not a string")

	// Validate expected content
	Expect(errorMsgStr).To(ContainSubstring(expectedContent))
}
