//nolint:dupl
package handlers_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/components/fun-app/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/noop"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func decodeStudentResponse(responseRecorder *httptest.ResponseRecorder) fun.Student {
	Expect(responseRecorder.Code).To(Equal(http.StatusCreated))

	var response fun.Student
	Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())

	return response
}

var _ = Describe("StudentHandler CUD", func() {
	var (
		ctx              context.Context
		db               *gorm.DB
		dbSQL            *sql.DB
		studentManager   manager.StudentManagerInterface
		router           *gin.Engine
		request          fun.StudentRequest
		updateRequest    fun.StudentRequest
		existingStudent  fun.Student
		response         fun.Student
		persisted        fun.Student
		persistedErr     common.HttpError
		audit            []fun.StudentAudit
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

		tracer := otel.Tracer("fun-app-student-handler-test")
		studentManager = manager.NewStudentManager(repository.NewStudentRepository(util.NewBaseDbRepository(db)), tracer)
		meter := noop.NewMeterProvider().Meter("fun-app-student-handler-test")
		createCounter, err := meter.Int64Counter("create_student")
		Expect(err).ToNot(HaveOccurred())
		studentCounter, err := meter.Int64UpDownCounter("student_count")
		Expect(err).ToNot(HaveOccurred())
		studentCreateTime, err := meter.Float64Histogram("student_create_time")
		Expect(err).ToNot(HaveOccurred())
		studentHandler := &handlers.StudentHandlerImpl{
			Manager:           studentManager,
			Tracer:            tracer,
			CreateCounter:     createCounter,
			StudentCounter:    studentCounter,
			StudentCreateTime: studentCreateTime,
		}

		router = util.CreateTestGinRouter()
		router.POST("/v1/student", studentHandler.CreateStudent)
		router.PUT("/v1/student/:id", studentHandler.UpdateStudent)
		router.DELETE("/v1/student/:id", studentHandler.DeleteStudents)
	})

	AfterEach(func() {
		if dbSQL != nil {
			Expect(dbSQL.Close()).To(Succeed())
		}
	})

	Context("POST /v1/student", func() {
		postStudent := func(studentRequest fun.StudentRequest) *httptest.ResponseRecorder {
			req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/student", studentRequest)
			router.ServeHTTP(recorder, req)
			return recorder
		}

		Context("Happy Path", func() {
			BeforeEach(func() {
				request = fun.StudentRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"}
				req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/student", request)
				responseRecorder = recorder
				router.ServeHTTP(responseRecorder, req)
				Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
				var err error
				persisted, err = studentManager.GetStudent(ctx, response.Id)
				Expect(err).ToNot(HaveOccurred())
				audit, err = studentManager.ListStudentAudit(ctx, response.Id)
				Expect(err).ToNot(HaveOccurred())
			})

			It("creates and persists the student with a CREATE audit", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusCreated))
				Expect(response.Id).ToNot(BeEmpty())
				Expect(response.Name).To(Equal(request.Name))
				Expect(response.Age).To(Equal(request.Age))
				Expect(response.Gender).To(Equal(request.Gender))
				Expect(persisted).To(Equal(response))
				Expect(audit).To(HaveLen(1))
				Expect(audit[0].Id).To(Equal(response.Id))
				Expect(audit[0].Name).To(Equal(request.Name))
				Expect(audit[0].Age).To(Equal(request.Age))
				Expect(audit[0].Gender).To(Equal(request.Gender))
				Expect(audit[0].Operation).To(Equal("CREATE"))
				Expect(audit[0].CreatedBy).To(Equal(fun.CreatedByAman))
				Expect(audit[0].CreatedAt).ToNot(BeZero())
				Expect(audit[0].CreatedAt).To(BeTemporally("~", time.Now(), 5*time.Second))
			})
		})

		Context("Field Validations", func() {
			Context("Name Field", func() {
				Context("Allowed Values", func() {
					Context("with the minimum length", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "A", Age: 36, Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("accepts a one-character Name", func() {
							response := decodeStudentResponse(responseRecorder)
							Expect(response.Name).To(Equal("A"))
						})
					})

					Context("with letters and spaces", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("accepts a Name with letters and spaces", func() {
							response := decodeStudentResponse(responseRecorder)
							Expect(response.Name).To(Equal("Ada Lovelace"))
						})
					})

					Context("with a hyphen", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Jean-Luc Picard", Age: 36, Gender: "MALE"}
							responseRecorder = postStudent(request)
						})

						It("accepts a hyphenated Name", func() {
							response := decodeStudentResponse(responseRecorder)
							Expect(response.Name).To(Equal("Jean-Luc Picard"))
						})
					})

					Context("with digits", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Agent 007", Age: 36, Gender: "MALE"}
							responseRecorder = postStudent(request)
						})

						It("accepts a Name with digits", func() {
							response := decodeStudentResponse(responseRecorder)
							Expect(response.Name).To(Equal("Agent 007"))
						})
					})

					Context("with the maximum length", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: strings.Repeat("A", 25), Age: 36, Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("accepts a 25-character Name", func() {
							response := decodeStudentResponse(responseRecorder)
							Expect(response.Name).To(Equal(strings.Repeat("A", 25)))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with a missing Name", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Age: 36, Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("returns a required Name validation error", func() {
							util.AssertError(responseRecorder, "Name", "required")
						})
					})

					Context("with an invalid character", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "A*B", Age: 36, Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("returns a Name character validation error", func() {
							util.AssertError(responseRecorder, "Name", "name")
						})
					})

					Context("with a 26-character Name", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: strings.Repeat("A", 26), Age: 36, Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("returns a maximum Name validation error", func() {
							util.AssertError(responseRecorder, "Name", "max")
						})
					})
				})
			})

			Context("Age Field", func() {
				Context("Allowed Values", func() {
					Context("with the minimum age", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Ada", Age: 1, Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("accepts age 1", func() {
							response := decodeStudentResponse(responseRecorder)
							Expect(response.Age).To(Equal(1))
						})
					})

					Context("with the maximum age", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Ada", Age: 150, Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("accepts age 150", func() {
							response := decodeStudentResponse(responseRecorder)
							Expect(response.Age).To(Equal(150))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with missing or zero age", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Ada", Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("returns a required Age validation error", func() {
							util.AssertError(responseRecorder, "Age", "required")
						})
					})

					Context("below the minimum age", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Ada", Age: -1, Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("returns a minimum Age validation error", func() {
							util.AssertError(responseRecorder, "Age", "min")
						})
					})

					Context("above the maximum age", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Ada", Age: 151, Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("returns a maximum Age validation error", func() {
							util.AssertError(responseRecorder, "Age", "max")
						})
					})
				})
			})

			Context("Gender Field", func() {
				Context("Allowed Values", func() {
					Context("with MALE", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Ada", Age: 36, Gender: "MALE"}
							responseRecorder = postStudent(request)
						})

						It("accepts MALE", func() {
							response := decodeStudentResponse(responseRecorder)
							Expect(response.Gender).To(Equal("MALE"))
						})
					})

					Context("with FEMALE", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Ada", Age: 36, Gender: "FEMALE"}
							responseRecorder = postStudent(request)
						})

						It("accepts FEMALE", func() {
							response := decodeStudentResponse(responseRecorder)
							Expect(response.Gender).To(Equal("FEMALE"))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with a missing Gender", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Ada", Age: 36}
							responseRecorder = postStudent(request)
						})

						It("returns a required Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "required")
						})
					})

					Context("with an unsupported Gender", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Ada", Age: 36, Gender: "OTHER"}
							responseRecorder = postStudent(request)
						})

						It("returns an equality Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "eq")
						})
					})

					Context("with lowercase gender", func() {
						BeforeEach(func() {
							request = fun.StudentRequest{Name: "Ada", Age: 36, Gender: "female"}
							responseRecorder = postStudent(request)
						})

						It("returns an equality Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "eq")
						})
					})
				})
			})
		})

		Context("Errors", func() {
			Context("malformed JSON", func() {
				BeforeEach(func() {
					req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/v1/student", strings.NewReader(`{"name":`))
					req.Header.Set("Content-Type", "application/json")
					responseRecorder = httptest.NewRecorder()
					router.ServeHTTP(responseRecorder, req)
				})

				It("returns HTTP 400", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
				})
			})
		})
	})

	Describe("PUT /v1/student/:id", func() {
		submitStudentUpdate := func(studentRequest fun.StudentRequest) {
			var err error
			existingStudent, err = studentManager.CreateStudent(ctx, fun.StudentRequest{
				Name:   "Ada Lovelace",
				Age:    36,
				Gender: "FEMALE",
			})
			Expect(err).ToNot(HaveOccurred())

			req, recorder := util.CreateTestRequest(http.MethodPut, "/v1/student/"+existingStudent.Id, studentRequest)
			responseRecorder = recorder
			router.ServeHTTP(responseRecorder, req)

			persisted, err = studentManager.GetStudent(ctx, existingStudent.Id)
			Expect(err).ToNot(HaveOccurred())
		}

		Context("Happy Path", func() {
			BeforeEach(func() {
				updateRequest = fun.StudentRequest{Name: "Grace Hopper", Age: 85, Gender: "FEMALE"}
				submitStudentUpdate(updateRequest)

				var err error
				audit, err = studentManager.ListStudentAudit(ctx, existingStudent.Id)
				Expect(err).ToNot(HaveOccurred())
			})

			It("updates and persists the student with CREATE then UPDATE audits", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
				Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
				Expect(persisted.Id).To(Equal(existingStudent.Id))
				Expect(persisted.Name).To(Equal(updateRequest.Name))
				Expect(persisted.Age).To(Equal(updateRequest.Age))
				Expect(persisted.Gender).To(Equal(updateRequest.Gender))

				Expect(audit).To(HaveLen(2))
				Expect(audit[0].Operation).To(Equal("CREATE"))
				Expect(audit[1].Id).To(Equal(existingStudent.Id))
				Expect(audit[1].Name).To(Equal(updateRequest.Name))
				Expect(audit[1].Age).To(Equal(updateRequest.Age))
				Expect(audit[1].Gender).To(Equal(updateRequest.Gender))
				Expect(audit[1].Operation).To(Equal("UPDATE"))
				Expect(audit[1].CreatedBy).To(Equal(fun.CreatedByAman))
				Expect(audit[1].CreatedAt).ToNot(BeZero())
			})
		})

		Context("Field Validations", func() {
			Context("Name Field", func() {
				Context("Allowed Values", func() {
					Context("with the minimum length", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "A", Age: 36, Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("accepts a one-character Name", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Name).To(Equal("A"))
						})
					})

					Context("with letters and spaces", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Grace Hopper", Age: 36, Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("accepts a Name with letters and spaces", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Name).To(Equal("Grace Hopper"))
						})
					})

					Context("with a hyphen", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Jean-Luc Picard", Age: 36, Gender: "MALE"}
							submitStudentUpdate(updateRequest)
						})

						It("accepts a hyphenated Name", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Name).To(Equal("Jean-Luc Picard"))
						})
					})

					Context("with digits", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Agent 007", Age: 36, Gender: "MALE"}
							submitStudentUpdate(updateRequest)
						})

						It("accepts a Name with digits", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Name).To(Equal("Agent 007"))
						})
					})

					Context("with the maximum length", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: strings.Repeat("A", 25), Age: 36, Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("accepts a 25-character Name", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Name).To(Equal(strings.Repeat("A", 25)))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with a missing Name", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Age: 36, Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("returns a required Name validation error", func() {
							util.AssertError(responseRecorder, "Name", "required")
						})
					})

					Context("with an invalid character", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "A*B", Age: 36, Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("returns a Name character validation error", func() {
							util.AssertError(responseRecorder, "Name", "name")
						})
					})

					Context("with a 26-character Name", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: strings.Repeat("A", 26), Age: 36, Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("returns a maximum Name validation error", func() {
							util.AssertError(responseRecorder, "Name", "max")
						})
					})
				})
			})

			Context("Age Field", func() {
				Context("Allowed Values", func() {
					Context("with the minimum age", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Ada", Age: 1, Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("accepts age 1", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Age).To(Equal(1))
						})
					})

					Context("with the maximum age", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Ada", Age: 150, Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("accepts age 150", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Age).To(Equal(150))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with missing or zero age", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Ada", Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("returns a required Age validation error", func() {
							util.AssertError(responseRecorder, "Age", "required")
						})
					})

					Context("below the minimum age", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Ada", Age: -1, Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("returns a minimum Age validation error", func() {
							util.AssertError(responseRecorder, "Age", "min")
						})
					})

					Context("above the maximum age", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Ada", Age: 151, Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("returns a maximum Age validation error", func() {
							util.AssertError(responseRecorder, "Age", "max")
						})
					})
				})
			})

			Context("Gender Field", func() {
				Context("Allowed Values", func() {
					Context("with MALE", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Ada", Age: 36, Gender: "MALE"}
							submitStudentUpdate(updateRequest)
						})

						It("accepts MALE", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Gender).To(Equal("MALE"))
						})
					})

					Context("with FEMALE", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Ada", Age: 36, Gender: "FEMALE"}
							submitStudentUpdate(updateRequest)
						})

						It("accepts FEMALE", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Gender).To(Equal("FEMALE"))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with a missing Gender", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Ada", Age: 36}
							submitStudentUpdate(updateRequest)
						})

						It("returns a required Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "required")
						})
					})

					Context("with an unsupported Gender", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Ada", Age: 36, Gender: "OTHER"}
							submitStudentUpdate(updateRequest)
						})

						It("returns an equality Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "eq")
						})
					})

					Context("with lowercase gender", func() {
						BeforeEach(func() {
							updateRequest = fun.StudentRequest{Name: "Ada", Age: 36, Gender: "female"}
							submitStudentUpdate(updateRequest)
						})

						It("returns an equality Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "eq")
						})
					})
				})
			})
		})

		Context("Errors", func() {
			Context("malformed JSON", func() {
				BeforeEach(func() {
					var err error
					existingStudent, err = studentManager.CreateStudent(ctx, fun.StudentRequest{
						Name:   "Ada Lovelace",
						Age:    36,
						Gender: "FEMALE",
					})
					Expect(err).ToNot(HaveOccurred())

					req := httptest.NewRequestWithContext(ctx, http.MethodPut, "/v1/student/"+existingStudent.Id, strings.NewReader(`{"name":`))
					req.Header.Set("Content-Type", "application/json")
					responseRecorder = httptest.NewRecorder()
					router.ServeHTTP(responseRecorder, req)
				})

				It("returns HTTP 400", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
				})
			})
		})
	})

	Describe("DELETE /v1/student/:id", func() {
		Context("Errors", func() {
			Context("with an empty ID", func() {
				BeforeEach(func() {
					req := httptest.NewRequestWithContext(ctx, http.MethodDelete, "/v1/student/", nil)
					responseRecorder = httptest.NewRecorder()
					router.ServeHTTP(responseRecorder, req)
				})

				It("returns HTTP 404", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
				})
			})

			Context("with a syntactically valid but absent ID", func() {
				BeforeEach(func() {
					req := httptest.NewRequestWithContext(ctx, http.MethodDelete, "/v1/student/missing-id", nil)
					responseRecorder = httptest.NewRecorder()
					router.ServeHTTP(responseRecorder, req)
				})

				It("returns HTTP 404", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
				})
			})
		})

		Context("Happy Path", func() {
			BeforeEach(func() {
				var err common.HttpError
				existingStudent, err = studentManager.CreateStudent(ctx, fun.StudentRequest{
					Name:   "Ada Lovelace",
					Age:    36,
					Gender: "FEMALE",
				})
				Expect(err).ToNot(HaveOccurred())

				req := httptest.NewRequestWithContext(ctx, http.MethodDelete, "/v1/student/"+existingStudent.Id, nil)
				responseRecorder = httptest.NewRecorder()
				router.ServeHTTP(responseRecorder, req)

				persisted, persistedErr = studentManager.GetStudent(ctx, existingStudent.Id)
				audit, err = studentManager.ListStudentAudit(ctx, existingStudent.Id)
				Expect(err).ToNot(HaveOccurred())
			})

			It("deletes the student with HTTP 204", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusNoContent))
				Expect(persistedErr).To(Equal(common.ErrNotFound))
			})

			It("records CREATE then DELETE audits", func() {
				Expect(audit).To(HaveLen(2))
				Expect(audit[0].Operation).To(Equal("CREATE"))
				Expect(audit[1].Id).To(Equal(existingStudent.Id))
				Expect(audit[1].Name).To(Equal(existingStudent.Name))
				Expect(audit[1].Age).To(Equal(existingStudent.Age))
				Expect(audit[1].Gender).To(Equal(existingStudent.Gender))
				Expect(audit[1].Operation).To(Equal("DELETE"))
				Expect(audit[1].CreatedBy).To(Equal(fun.CreatedByAman))
				Expect(audit[1].CreatedAt).ToNot(BeZero())
			})
		})
	})
})
