package main

// TestMain - test to drive external testing coverage
/**
go test fun-app_test.go fun-app.go -coverprofile=coverage.out
curl http://localhost:8080/admin/stop
go tool cover -func=coverage.out
*/
//func TestFunApp(t *testing.T) {
//	main()
//TODO: Solve for coverage.