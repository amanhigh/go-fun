//nolint:dupl
package handlers_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/components/fun-app/repository"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/noop"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type studentGetEnrollmentHandlerStub struct{}

func (studentGetEnrollmentHandlerStub) CreateEnrollment(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
func (studentGetEnrollmentHandlerStub) GetEnrollment(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

type studentGetAdminHandlerStub struct{}

func (studentGetAdminHandlerStub) Stop(c *gin.Context) { c.Status(http.StatusNotImplemented) }

func decodeStudentGetResponse(responseRecorder *httptest.ResponseRecorder) fun.Student {
	var response fun.Student
	Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
	return response
}

func decodeStudentListResponse(responseRecorder *httptest.ResponseRecorder) fun.StudentList {
	var response fun.StudentList
	Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
	return response
}

func decodeStudentAuditResponse(responseRecorder *httptest.ResponseRecorder) []fun.StudentAudit {
	var response []fun.StudentAudit
	Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
	return response
}

func decodeStudentGetError(responseRecorder *httptest.ResponseRecorder) map[string]any {
	var response map[string]any
	Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
	return response
}

var _ = Describe("Student Handler Integration - GET Tests", func() {
	var (
		ctx              context.Context
		db               *gorm.DB
		dbSQL            *sql.DB
		studentManager   manager.StudentManagerInterface
		router           *gin.Engine
		existingStudent  fun.Student
		responseRecorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		db, err = util.CreateTestDb(gormlogger.Warn)
		Expect(err).ToNot(HaveOccurred())
		Expect(db.AutoMigrate(&fun.Student{}, &fun.StudentAudit{})).To(Succeed())
		dbSQL, err = db.DB()
		Expect(err).ToNot(HaveOccurred())

		tracer := otel.Tracer("fun-app-student-get-handler-test")
		studentManager = manager.NewStudentManager(repository.NewStudentRepository(util.NewBaseDbRepository(db)), tracer)
		meter := noop.NewMeterProvider().Meter("fun-app-student-get-handler-test")
		createCounter, err := meter.Int64Counter("get_test_create_student")
		Expect(err).ToNot(HaveOccurred())
		studentCounter, err := meter.Int64UpDownCounter("get_test_student_count")
		Expect(err).ToNot(HaveOccurred())
		studentCreateTime, err := meter.Float64Histogram("get_test_student_create_time")
		Expect(err).ToNot(HaveOccurred())
		studentHandler := &handlers.StudentHandlerImpl{
			Manager:           studentManager,
			Tracer:            tracer,
			CreateCounter:     createCounter,
			StudentCounter:    studentCounter,
			StudentCreateTime: studentCreateTime,
		}

		lifecycle := &handlers.FunAppServerLifecycle{
			StudentHandler:    studentHandler,
			EnrollmentHandler: studentGetEnrollmentHandlerStub{},
			AdminHandler:      studentGetAdminHandlerStub{},
		}
		router = util.CreateTestGinRouter()
		lifecycle.RegisterRoutes(router)
	})

	AfterEach(func() {
		if dbSQL != nil {
			Expect(dbSQL.Close()).To(Succeed())
		}
	})

	Describe("GET /v1/student/:id", func() {
		Context("Happy Path", func() {
			BeforeEach(func() {
				var err error
				existingStudent, err = studentManager.CreateStudent(ctx, fun.StudentRequest{
					Name: "Ada Lovelace", Age: 36, Gender: "FEMALE",
				})
				Expect(err).ToNot(HaveOccurred())
				req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/student/"+existingStudent.Id, nil)
				responseRecorder = httptest.NewRecorder()
				router.ServeHTTP(responseRecorder, req)
			})

			It("returns all persisted student fields", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
				response := decodeStudentGetResponse(responseRecorder)
				Expect(response.Id).To(Equal(existingStudent.Id))
				Expect(response.Name).To(Equal(existingStudent.Name))
				Expect(response.Age).To(Equal(existingStudent.Age))
				Expect(response.Gender).To(Equal(existingStudent.Gender))
			})
		})

		Context("Field Validations", func() {
			Context("Student ID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						var err error
						existingStudent, err = studentManager.CreateStudent(ctx, fun.StudentRequest{Name: "Ada", Age: 36, Gender: "FEMALE"})
						Expect(err).ToNot(HaveOccurred())
						req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/student/"+existingStudent.Id, nil)
						responseRecorder = httptest.NewRecorder()
						router.ServeHTTP(responseRecorder, req)
					})

					It("accepts a created and persisted student ID", func() {
						Expect(responseRecorder.Code).To(Equal(http.StatusOK))
						response := decodeStudentGetResponse(responseRecorder)
						Expect(response.Id).To(Equal(existingStudent.Id))
					})
				})

				Context("Bad Values", func() {
					BeforeEach(func() {
						req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/student/missing-student", nil)
						responseRecorder = httptest.NewRecorder()
						router.ServeHTTP(responseRecorder, req)
					})

					It("returns a meaningful JSend not-found failure", func() {
						Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
						response := decodeStudentGetError(responseRecorder)
						Expect(response["status"]).To(Equal("fail"))
						data, ok := response["data"].(map[string]any)
						Expect(ok).To(BeTrue())
						Expect(data["message"]).To(Equal("NotFound"))
					})
				})
			})
		})

		Context("Errors", func() {
			BeforeEach(func() {
				Expect(db.Migrator().DropTable(&fun.Student{})).To(Succeed())
				req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/student/fixed-id", nil)
				responseRecorder = httptest.NewRecorder()
				router.ServeHTTP(responseRecorder, req)
			})

			It("returns HTTP 500 for an unavailable student table", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
				response := decodeStudentGetError(responseRecorder)
				Expect(response["status"]).To(Equal("error"))
				Expect(response["message"]).ToNot(BeEmpty())
				Expect(response["code"]).To(Equal(float64(http.StatusInternalServerError)))
			})
		})
	})

	Describe("GET /v1/student/:id/audit", func() {
		Context("Happy Path", func() {
			var audit []fun.StudentAudit

			BeforeEach(func() {
				var err error
				existingStudent, err = studentManager.CreateStudent(ctx, fun.StudentRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"})
				Expect(err).ToNot(HaveOccurred())
				Expect(studentManager.UpdateStudent(ctx, existingStudent.Id, fun.StudentRequest{Name: "Grace Hopper", Age: 85, Gender: "FEMALE"})).ToNot(HaveOccurred())
				req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/student/"+existingStudent.Id+"/audit", nil)
				responseRecorder = httptest.NewRecorder()
				router.ServeHTTP(responseRecorder, req)
				audit = decodeStudentAuditResponse(responseRecorder)
			})

			It("returns ordered CREATE and UPDATE audit records", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
				Expect(audit).To(HaveLen(2))
				Expect(audit[0].AuditID).ToNot(BeZero())
				Expect(audit[0].Id).To(Equal(existingStudent.Id))
				Expect(audit[0].Name).To(Equal("Ada Lovelace"))
				Expect(audit[0].Age).To(Equal(36))
				Expect(audit[0].Gender).To(Equal("FEMALE"))
				Expect(audit[0].Operation).To(Equal("CREATE"))
				Expect(audit[0].CreatedBy).To(Equal(fun.CreatedByAman))
				Expect(audit[0].CreatedAt).ToNot(BeZero())
				Expect(audit[1].AuditID).ToNot(BeZero())
				Expect(audit[1].Id).To(Equal(existingStudent.Id))
				Expect(audit[1].Name).To(Equal("Grace Hopper"))
				Expect(audit[1].Age).To(Equal(85))
				Expect(audit[1].Gender).To(Equal("FEMALE"))
				Expect(audit[1].Operation).To(Equal("UPDATE"))
				Expect(audit[1].CreatedBy).To(Equal(fun.CreatedByAman))
				Expect(audit[1].CreatedAt).ToNot(BeZero())
				Expect(audit[1].CreatedAt).To(BeTemporally(">=", audit[0].CreatedAt))
			})
		})

		Context("Errors", func() {
			BeforeEach(func() {
				Expect(db.Migrator().DropTable(&fun.StudentAudit{})).To(Succeed())
				req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/student/missing-student/audit", nil)
				responseRecorder = httptest.NewRecorder()
				router.ServeHTTP(responseRecorder, req)
			})

			It("returns HTTP 500 for an unavailable audit table", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
				response := decodeStudentGetError(responseRecorder)
				Expect(response["status"]).To(Equal("error"))
				Expect(response["message"]).ToNot(BeEmpty())
				Expect(response["code"]).To(Equal(float64(http.StatusInternalServerError)))
			})
		})
	})

	Describe("GET /v1/student/", func() {
		createStudents := func(requests ...fun.StudentRequest) {
			for _, request := range requests {
				_, err := studentManager.CreateStudent(ctx, request)
				Expect(err).ToNot(HaveOccurred())
			}
		}
		listStudentsResponseOnly := func(query string) *httptest.ResponseRecorder {
			req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/student/"+query, nil)
			responseRecorder = httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			return responseRecorder
		}
		listStudents := func(query string) fun.StudentList {
			return decodeStudentListResponse(listStudentsResponseOnly(query))
		}

		Context("Happy Path", func() {
			Context("with no optional query fields", func() {
				var response fun.StudentList

				BeforeEach(func() {
					requests := make([]fun.StudentRequest, 22)
					for i := range requests {
						// Descending name order: Student V, Student U, ..., Student A
						requests[i] = fun.StudentRequest{Name: "Student " + string(rune('V'-i)), Age: 22 - i, Gender: "FEMALE"}
					}
					createStudents(requests...)
					response = listStudents("")
				})

				It("returns a raw list with default pagination and total metadata", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusOK))
					Expect(response.Records).To(HaveLen(20))
					Expect(response.Metadata.Offset).To(Equal(0))
					Expect(response.Metadata.Limit).To(Equal(20))
					Expect(response.Metadata.Total).To(Equal(int64(22)))
					Expect(response.Records[0].Name).To(Equal("Student A"))
					Expect(response.Records[19].Name).To(Equal("Student T"))
				})
			})

			Context("with combined name and gender filters proving AND semantics with empty result", func() {
				var response fun.StudentList

				BeforeEach(func() {
					createStudents(fun.StudentRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"})
					response = listStudents("?name=Ada&gender=MALE")
				})

				It("returns an empty list and zero total", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusOK))
					Expect(response.Records).To(BeEmpty())
					Expect(response.Metadata.Total).To(Equal(int64(0)))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Offset Field", func() {
				Context("Allowed Values", func() {
					Context("with offset 0 and limit 2", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Bob", Age: 40, Gender: "MALE"},
								fun.StudentRequest{Name: "Carol", Age: 28, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Dave", Age: 50, Gender: "MALE"},
							)
							response = listStudents("?offset=0&limit=2&sort_by=name&sort-order=asc")
						})

						It("returns the first page of records with correct metadata", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(2))
							Expect(response.Records[0].Name).To(Equal("Ada"))
							Expect(response.Records[1].Name).To(Equal("Bob"))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(2))
							Expect(response.Metadata.Total).To(Equal(int64(4)))
						})
					})

					Context("with a positive offset and limit 2", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Bob", Age: 40, Gender: "MALE"},
								fun.StudentRequest{Name: "Carol", Age: 28, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Dave", Age: 50, Gender: "MALE"},
							)
							response = listStudents("?offset=2&limit=2&sort_by=name&sort-order=asc")
						})

						It("returns the second page of records with correct metadata", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(2))
							Expect(response.Records[0].Name).To(Equal("Carol"))
							Expect(response.Records[1].Name).To(Equal("Dave"))
							Expect(response.Metadata.Offset).To(Equal(2))
							Expect(response.Metadata.Limit).To(Equal(2))
							Expect(response.Metadata.Total).To(Equal(int64(4)))
						})
					})
				})

				Context("Bad Values", func() {
					BeforeEach(func() {
						responseRecorder = listStudentsResponseOnly("?offset=-1")
					})

					It("returns a minimum Offset validation error", func() {
						util.AssertError(responseRecorder, "Offset", "min")
					})
				})
			})

			Context("Limit Field", func() {
				Context("Allowed Values", func() {
					Context("with limit 1", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Bob", Age: 40, Gender: "MALE"},
								fun.StudentRequest{Name: "Carol", Age: 28, Gender: "FEMALE"},
							)
							response = listStudents("?limit=1&sort_by=name&sort-order=asc")
						})

						It("returns one record with correct metadata", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(1))
							Expect(response.Records[0].Name).To(Equal("Ada"))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(1))
							Expect(response.Metadata.Total).To(Equal(int64(3)))
						})
					})

					Context("with limit 100", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Bob", Age: 40, Gender: "MALE"},
								fun.StudentRequest{Name: "Carol", Age: 28, Gender: "FEMALE"},
							)
							response = listStudents("?limit=100&sort_by=name&sort-order=asc")
						})

						It("returns all records with correct metadata", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(3))
							Expect(response.Records[0].Name).To(Equal("Ada"))
							Expect(response.Records[1].Name).To(Equal("Bob"))
							Expect(response.Records[2].Name).To(Equal("Carol"))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(100))
							Expect(response.Metadata.Total).To(Equal(int64(3)))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with limit 0", func() {
						BeforeEach(func() { responseRecorder = listStudentsResponseOnly("?limit=0") })

						It("returns a minimum Limit validation error", func() {
							util.AssertError(responseRecorder, "Limit", "min")
						})
					})

					Context("with limit 101", func() {
						BeforeEach(func() { responseRecorder = listStudentsResponseOnly("?limit=101") })

						It("returns a maximum Limit validation error", func() {
							util.AssertError(responseRecorder, "Limit", "max")
						})
					})
				})
			})

			Context("Name Field", func() {
				Context("Allowed Values", func() {
					Context("with a partial name", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Ada Byron", Age: 28, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Grace Hopper", Age: 85, Gender: "FEMALE"},
							)
							response = listStudents("?name=Ada&sort_by=name")
						})

						It("returns matching records and total metadata", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(2))
							Expect(response.Records[0].Name).To(Equal("Ada Byron"))
							Expect(response.Records[1].Name).To(Equal("Ada Lovelace"))
							Expect(response.Metadata.Total).To(Equal(int64(2)))
						})
					})

					Context("with a 25-character name", func() {
						var response fun.StudentList

						BeforeEach(func() {
							name25 := strings.Repeat("A", 25)
							createStudents(
								fun.StudentRequest{Name: name25, Age: 30, Gender: "MALE"},
								fun.StudentRequest{Name: "Bob", Age: 40, Gender: "MALE"},
							)
							response = listStudents("?name=" + name25)
						})

						It("returns the matching record", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(1))
							Expect(response.Records[0].Name).To(Equal(strings.Repeat("A", 25)))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with an invalid character", func() {
						BeforeEach(func() { responseRecorder = listStudentsResponseOnly("?name=A%2AB") })

						It("returns a Name character validation error", func() {
							util.AssertError(responseRecorder, "Name", "name")
						})
					})

					Context("with a 26-character name", func() {
						BeforeEach(func() {
							responseRecorder = listStudentsResponseOnly("?name=" + strings.Repeat("A", 26))
						})

						It("returns a maximum Name validation error", func() {
							util.AssertError(responseRecorder, "Name", "max")
						})
					})
				})
			})

			Context("Gender Field", func() {
				Context("Allowed Values", func() {
					Context("with MALE", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Alan", Age: 42, Gender: "MALE"},
								fun.StudentRequest{Name: "Bob", Age: 30, Gender: "MALE"},
								fun.StudentRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
							)
							response = listStudents("?gender=MALE")
						})

						It("returns only MALE records with correct total", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(2))
							Expect(response.Records[0].Gender).To(Equal("MALE"))
							Expect(response.Records[1].Gender).To(Equal("MALE"))
							Expect(response.Metadata.Total).To(Equal(int64(2)))
						})
					})

					Context("with FEMALE", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Grace", Age: 85, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Bob", Age: 30, Gender: "MALE"},
							)
							response = listStudents("?gender=FEMALE")
						})

						It("returns only FEMALE records with correct total", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(2))
							Expect(response.Records[0].Gender).To(Equal("FEMALE"))
							Expect(response.Records[1].Gender).To(Equal("FEMALE"))
							Expect(response.Metadata.Total).To(Equal(int64(2)))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with unsupported OTHER", func() {
						BeforeEach(func() { responseRecorder = listStudentsResponseOnly("?gender=OTHER") })

						It("returns an equality Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "eq")
						})
					})

					Context("with lowercase female", func() {
						BeforeEach(func() { responseRecorder = listStudentsResponseOnly("?gender=female") })

						It("returns an equality Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "eq")
						})
					})
				})
			})

			Context("SortBy Field", func() {
				Context("Allowed Values", func() {
					Context("with name", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Charlie", Age: 30, Gender: "MALE"},
								fun.StudentRequest{Name: "Alice", Age: 25, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Bob", Age: 40, Gender: "MALE"},
							)
							response = listStudents("?sort_by=name")
						})

						It("returns records in ascending name order", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(3))
							Expect(response.Records[0].Name).To(Equal("Alice"))
							Expect(response.Records[1].Name).To(Equal("Bob"))
							Expect(response.Records[2].Name).To(Equal("Charlie"))
						})
					})

					Context("with age", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Older", Age: 42, Gender: "MALE"},
								fun.StudentRequest{Name: "Youngest", Age: 18, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Middle", Age: 30, Gender: "FEMALE"},
							)
							response = listStudents("?sort_by=age")
						})

						It("returns records in ascending age order", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(3))
							Expect(response.Records[0].Name).To(Equal("Youngest"))
							Expect(response.Records[1].Name).To(Equal("Middle"))
							Expect(response.Records[2].Name).To(Equal("Older"))
						})
					})

					Context("with gender", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Male One", Age: 30, Gender: "MALE"},
								fun.StudentRequest{Name: "Female One", Age: 25, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Male Two", Age: 40, Gender: "MALE"},
							)
							response = listStudents("?sort_by=gender")
						})

						It("returns records in ascending gender groups", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(3))
							Expect(response.Records[0].Gender).To(Equal("FEMALE"))
							Expect(response.Records[1].Gender).To(Equal("MALE"))
							Expect(response.Records[2].Gender).To(Equal("MALE"))
						})
					})
				})

				Context("Bad Values", func() {
					BeforeEach(func() { responseRecorder = listStudentsResponseOnly("?sort_by=invalid") })

					It("returns an equality SortBy validation error", func() {
						util.AssertError(responseRecorder, "SortBy", "eq")
					})
				})
			})

			Context("SortOrder Field", func() {
				Context("Allowed Values", func() {
					Context("with asc", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Charlie", Age: 30, Gender: "MALE"},
								fun.StudentRequest{Name: "Alice", Age: 25, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Bob", Age: 40, Gender: "MALE"},
							)
							response = listStudents("?sort_by=name&sort-order=asc")
						})

						It("returns records in ascending name order", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(3))
							Expect(response.Records[0].Name).To(Equal("Alice"))
							Expect(response.Records[1].Name).To(Equal("Bob"))
							Expect(response.Records[2].Name).To(Equal("Charlie"))
						})
					})

					Context("with desc", func() {
						var response fun.StudentList

						BeforeEach(func() {
							createStudents(
								fun.StudentRequest{Name: "Charlie", Age: 30, Gender: "MALE"},
								fun.StudentRequest{Name: "Alice", Age: 25, Gender: "FEMALE"},
								fun.StudentRequest{Name: "Bob", Age: 40, Gender: "MALE"},
							)
							response = listStudents("?sort_by=name&sort-order=desc")
						})

						It("returns records in descending name order", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Records).To(HaveLen(3))
							Expect(response.Records[0].Name).To(Equal("Charlie"))
							Expect(response.Records[1].Name).To(Equal("Bob"))
							Expect(response.Records[2].Name).To(Equal("Alice"))
						})
					})
				})

				Context("Bad Values", func() {
					BeforeEach(func() { responseRecorder = listStudentsResponseOnly("?sort-order=invalid") })

					It("returns a oneof SortOrder validation error", func() {
						util.AssertError(responseRecorder, "SortOrder", "oneof")
					})
				})
			})
		})

		Context("Errors", func() {
			BeforeEach(func() {
				Expect(db.Migrator().DropTable(&fun.Student{})).To(Succeed())
				responseRecorder = listStudentsResponseOnly("")
			})

			// Intentionally assert only status code: the raw response body format
			// (raw string vs JSend envelope) is an implementation detail not yet stabilised.
			It("returns HTTP 500 when the database table is missing", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
