package main

import (
	"testing"
)

// TestMain - test to drive external testing coverage
/**
go test app_test.go fun-app.go -coverprofile=coverage.out
go tool cover -func=coverage.out
*/
func TestFunApp(t *testing.T) {
	main()
}
