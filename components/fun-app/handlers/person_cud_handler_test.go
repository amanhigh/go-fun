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
	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
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

func decodePersonResponse(responseRecorder *httptest.ResponseRecorder) fun.Person {
	Expect(responseRecorder.Code).To(Equal(http.StatusCreated))

	var response fun.Person
	Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())

	return response
}

var _ = Describe("PersonHandler CUD", func() {
	var (
		ctx              context.Context
		db               *gorm.DB
		dbSQL            *sql.DB
		personManager    manager.PersonManagerInterface
		router           *gin.Engine
		request          fun.PersonRequest
		updateRequest    fun.PersonRequest
		existingPerson   fun.Person
		response         fun.Person
		persisted        fun.Person
		persistedErr     common.HttpError
		audit            []fun.PersonAudit
		responseRecorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		db, err = util.CreateTestDb(gormlogger.Warn)
		Expect(err).ToNot(HaveOccurred())
		Expect(db.AutoMigrate(&fun.Person{}, &fun.PersonAudit{})).To(Succeed())
		dbSQL, err = db.DB()
		Expect(err).ToNot(HaveOccurred())

		tracer := otel.Tracer("fun-app-person-handler-test")
		personManager = manager.NewPersonManager(dao.NewPersonDao(util.NewBaseDbRepository(db)), tracer)
		meter := noop.NewMeterProvider().Meter("fun-app-person-handler-test")
		createCounter, err := meter.Int64Counter("create_person")
		Expect(err).ToNot(HaveOccurred())
		personCounter, err := meter.Int64UpDownCounter("person_count")
		Expect(err).ToNot(HaveOccurred())
		personCreateTime, err := meter.Float64Histogram("person_create_time")
		Expect(err).ToNot(HaveOccurred())
		personHandler := &handlers.PersonHandlerImpl{
			Manager:          personManager,
			Tracer:           tracer,
			CreateCounter:    createCounter,
			PersonCounter:    personCounter,
			PersonCreateTime: personCreateTime,
		}

		router = util.CreateTestGinRouter()
		router.POST("/v1/person", personHandler.CreatePerson)
		router.PUT("/v1/person/:id", personHandler.UpdatePerson)
		router.DELETE("/v1/person/:id", personHandler.DeletePersons)
	})

	AfterEach(func() {
		if dbSQL != nil {
			Expect(dbSQL.Close()).To(Succeed())
		}
	})

	Context("POST /v1/person", func() {
		postPerson := func(personRequest fun.PersonRequest) *httptest.ResponseRecorder {
			req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/person", personRequest)
			router.ServeHTTP(recorder, req)
			return recorder
		}

		Context("Happy Path", func() {
			BeforeEach(func() {
				request = fun.PersonRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"}
				req, recorder := util.CreateTestRequest(http.MethodPost, "/v1/person", request)
				responseRecorder = recorder
				router.ServeHTTP(responseRecorder, req)
				Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
				var err error
				persisted, err = personManager.GetPerson(ctx, response.Id)
				Expect(err).ToNot(HaveOccurred())
				audit, err = personManager.ListPersonAudit(ctx, response.Id)
				Expect(err).ToNot(HaveOccurred())
			})

			It("creates and persists the person with a CREATE audit", func() {
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
							request = fun.PersonRequest{Name: "A", Age: 36, Gender: "FEMALE"}
							responseRecorder = postPerson(request)
						})

						It("accepts a one-character Name", func() {
							response := decodePersonResponse(responseRecorder)
							Expect(response.Name).To(Equal("A"))
						})
					})

					Context("with letters and spaces", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"}
							responseRecorder = postPerson(request)
						})

						It("accepts a Name with letters and spaces", func() {
							response := decodePersonResponse(responseRecorder)
							Expect(response.Name).To(Equal("Ada Lovelace"))
						})
					})

					Context("with a hyphen", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "Jean-Luc Picard", Age: 36, Gender: "MALE"}
							responseRecorder = postPerson(request)
						})

						It("accepts a hyphenated Name", func() {
							response := decodePersonResponse(responseRecorder)
							Expect(response.Name).To(Equal("Jean-Luc Picard"))
						})
					})

					Context("with digits", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "Agent 007", Age: 36, Gender: "MALE"}
							responseRecorder = postPerson(request)
						})

						It("accepts a Name with digits", func() {
							response := decodePersonResponse(responseRecorder)
							Expect(response.Name).To(Equal("Agent 007"))
						})
					})

					Context("with the maximum length", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: strings.Repeat("A", 25), Age: 36, Gender: "FEMALE"}
							responseRecorder = postPerson(request)
						})

						It("accepts a 25-character Name", func() {
							response := decodePersonResponse(responseRecorder)
							Expect(response.Name).To(Equal(strings.Repeat("A", 25)))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with a missing Name", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Age: 36, Gender: "FEMALE"}
							responseRecorder = postPerson(request)
						})

						It("returns a required Name validation error", func() {
							util.AssertError(responseRecorder, "Name", "required")
						})
					})

					Context("with an invalid character", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "A*B", Age: 36, Gender: "FEMALE"}
							responseRecorder = postPerson(request)
						})

						It("returns a Name character validation error", func() {
							util.AssertError(responseRecorder, "Name", "name")
						})
					})

					Context("with a 26-character Name", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: strings.Repeat("A", 26), Age: 36, Gender: "FEMALE"}
							responseRecorder = postPerson(request)
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
							request = fun.PersonRequest{Name: "Ada", Age: 1, Gender: "FEMALE"}
							responseRecorder = postPerson(request)
						})

						It("accepts age 1", func() {
							response := decodePersonResponse(responseRecorder)
							Expect(response.Age).To(Equal(1))
						})
					})

					Context("with the maximum age", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "Ada", Age: 150, Gender: "FEMALE"}
							responseRecorder = postPerson(request)
						})

						It("accepts age 150", func() {
							response := decodePersonResponse(responseRecorder)
							Expect(response.Age).To(Equal(150))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with missing or zero age", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "Ada", Gender: "FEMALE"}
							responseRecorder = postPerson(request)
						})

						It("returns a required Age validation error", func() {
							util.AssertError(responseRecorder, "Age", "required")
						})
					})

					Context("below the minimum age", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "Ada", Age: -1, Gender: "FEMALE"}
							responseRecorder = postPerson(request)
						})

						It("returns a minimum Age validation error", func() {
							util.AssertError(responseRecorder, "Age", "min")
						})
					})

					Context("above the maximum age", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "Ada", Age: 151, Gender: "FEMALE"}
							responseRecorder = postPerson(request)
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
							request = fun.PersonRequest{Name: "Ada", Age: 36, Gender: "MALE"}
							responseRecorder = postPerson(request)
						})

						It("accepts MALE", func() {
							response := decodePersonResponse(responseRecorder)
							Expect(response.Gender).To(Equal("MALE"))
						})
					})

					Context("with FEMALE", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"}
							responseRecorder = postPerson(request)
						})

						It("accepts FEMALE", func() {
							response := decodePersonResponse(responseRecorder)
							Expect(response.Gender).To(Equal("FEMALE"))
						})
					})
				})

				Context("Bad Values", func() {
					Context("with a missing Gender", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "Ada", Age: 36}
							responseRecorder = postPerson(request)
						})

						It("returns a required Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "required")
						})
					})

					Context("with an unsupported Gender", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "Ada", Age: 36, Gender: "OTHER"}
							responseRecorder = postPerson(request)
						})

						It("returns an equality Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "eq")
						})
					})

					Context("with lowercase gender", func() {
						BeforeEach(func() {
							request = fun.PersonRequest{Name: "Ada", Age: 36, Gender: "female"}
							responseRecorder = postPerson(request)
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
					req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/v1/person", strings.NewReader(`{"name":`))
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

	Describe("PUT /v1/person/:id", func() {
		submitPersonUpdate := func(personRequest fun.PersonRequest) {
			var err error
			existingPerson, err = personManager.CreatePerson(ctx, fun.PersonRequest{
				Name:   "Ada Lovelace",
				Age:    36,
				Gender: "FEMALE",
			})
			Expect(err).ToNot(HaveOccurred())

			req, recorder := util.CreateTestRequest(http.MethodPut, "/v1/person/"+existingPerson.Id, personRequest)
			responseRecorder = recorder
			router.ServeHTTP(responseRecorder, req)

			persisted, err = personManager.GetPerson(ctx, existingPerson.Id)
			Expect(err).ToNot(HaveOccurred())
		}

		Context("Happy Path", func() {
			BeforeEach(func() {
				updateRequest = fun.PersonRequest{Name: "Grace Hopper", Age: 85, Gender: "FEMALE"}
				submitPersonUpdate(updateRequest)

				var err error
				audit, err = personManager.ListPersonAudit(ctx, existingPerson.Id)
				Expect(err).ToNot(HaveOccurred())
			})

			It("updates and persists the person with CREATE then UPDATE audits", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
				Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
				Expect(persisted.Id).To(Equal(existingPerson.Id))
				Expect(persisted.Name).To(Equal(updateRequest.Name))
				Expect(persisted.Age).To(Equal(updateRequest.Age))
				Expect(persisted.Gender).To(Equal(updateRequest.Gender))

				Expect(audit).To(HaveLen(2))
				Expect(audit[0].Operation).To(Equal("CREATE"))
				Expect(audit[1].Id).To(Equal(existingPerson.Id))
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
							updateRequest = fun.PersonRequest{Name: "A", Age: 36, Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
						})

						It("accepts a one-character Name", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Name).To(Equal("A"))
						})
					})

					Context("with letters and spaces", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: "Grace Hopper", Age: 36, Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
						})

						It("accepts a Name with letters and spaces", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Name).To(Equal("Grace Hopper"))
						})
					})

					Context("with a hyphen", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: "Jean-Luc Picard", Age: 36, Gender: "MALE"}
							submitPersonUpdate(updateRequest)
						})

						It("accepts a hyphenated Name", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Name).To(Equal("Jean-Luc Picard"))
						})
					})

					Context("with digits", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: "Agent 007", Age: 36, Gender: "MALE"}
							submitPersonUpdate(updateRequest)
						})

						It("accepts a Name with digits", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Name).To(Equal("Agent 007"))
						})
					})

					Context("with the maximum length", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: strings.Repeat("A", 25), Age: 36, Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
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
							updateRequest = fun.PersonRequest{Age: 36, Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
						})

						It("returns a required Name validation error", func() {
							util.AssertError(responseRecorder, "Name", "required")
						})
					})

					Context("with an invalid character", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: "A*B", Age: 36, Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
						})

						It("returns a Name character validation error", func() {
							util.AssertError(responseRecorder, "Name", "name")
						})
					})

					Context("with a 26-character Name", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: strings.Repeat("A", 26), Age: 36, Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
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
							updateRequest = fun.PersonRequest{Name: "Ada", Age: 1, Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
						})

						It("accepts age 1", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Age).To(Equal(1))
						})
					})

					Context("with the maximum age", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: "Ada", Age: 150, Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
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
							updateRequest = fun.PersonRequest{Name: "Ada", Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
						})

						It("returns a required Age validation error", func() {
							util.AssertError(responseRecorder, "Age", "required")
						})
					})

					Context("below the minimum age", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: "Ada", Age: -1, Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
						})

						It("returns a minimum Age validation error", func() {
							util.AssertError(responseRecorder, "Age", "min")
						})
					})

					Context("above the maximum age", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: "Ada", Age: 151, Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
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
							updateRequest = fun.PersonRequest{Name: "Ada", Age: 36, Gender: "MALE"}
							submitPersonUpdate(updateRequest)
						})

						It("accepts MALE", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(responseRecorder.Body.String()).To(Equal(`"UPDATED"`))
							Expect(persisted.Gender).To(Equal("MALE"))
						})
					})

					Context("with FEMALE", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"}
							submitPersonUpdate(updateRequest)
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
							updateRequest = fun.PersonRequest{Name: "Ada", Age: 36}
							submitPersonUpdate(updateRequest)
						})

						It("returns a required Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "required")
						})
					})

					Context("with an unsupported Gender", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: "Ada", Age: 36, Gender: "OTHER"}
							submitPersonUpdate(updateRequest)
						})

						It("returns an equality Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "eq")
						})
					})

					Context("with lowercase gender", func() {
						BeforeEach(func() {
							updateRequest = fun.PersonRequest{Name: "Ada", Age: 36, Gender: "female"}
							submitPersonUpdate(updateRequest)
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
					existingPerson, err = personManager.CreatePerson(ctx, fun.PersonRequest{
						Name:   "Ada Lovelace",
						Age:    36,
						Gender: "FEMALE",
					})
					Expect(err).ToNot(HaveOccurred())

					req := httptest.NewRequestWithContext(ctx, http.MethodPut, "/v1/person/"+existingPerson.Id, strings.NewReader(`{"name":`))
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

	Describe("DELETE /v1/person/:id", func() {
		Context("Errors", func() {
			Context("with an empty ID", func() {
				BeforeEach(func() {
					req := httptest.NewRequestWithContext(ctx, http.MethodDelete, "/v1/person/", nil)
					responseRecorder = httptest.NewRecorder()
					router.ServeHTTP(responseRecorder, req)
				})

				It("returns HTTP 404", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
				})
			})

			Context("with a syntactically valid but absent ID", func() {
				BeforeEach(func() {
					req := httptest.NewRequestWithContext(ctx, http.MethodDelete, "/v1/person/missing-id", nil)
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
				existingPerson, err = personManager.CreatePerson(ctx, fun.PersonRequest{
					Name:   "Ada Lovelace",
					Age:    36,
					Gender: "FEMALE",
				})
				Expect(err).ToNot(HaveOccurred())

				req := httptest.NewRequestWithContext(ctx, http.MethodDelete, "/v1/person/"+existingPerson.Id, nil)
				responseRecorder = httptest.NewRecorder()
				router.ServeHTTP(responseRecorder, req)

				persisted, persistedErr = personManager.GetPerson(ctx, existingPerson.Id)
				audit, err = personManager.ListPersonAudit(ctx, existingPerson.Id)
				Expect(err).ToNot(HaveOccurred())
			})

			It("deletes the person with HTTP 204", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusNoContent))
				Expect(persistedErr).To(Equal(common.ErrNotFound))
			})

			It("records CREATE then DELETE audits", func() {
				Expect(audit).To(HaveLen(2))
				Expect(audit[0].Operation).To(Equal("CREATE"))
				Expect(audit[1].Id).To(Equal(existingPerson.Id))
				Expect(audit[1].Name).To(Equal(existingPerson.Name))
				Expect(audit[1].Age).To(Equal(existingPerson.Age))
				Expect(audit[1].Gender).To(Equal(existingPerson.Gender))
				Expect(audit[1].Operation).To(Equal("DELETE"))
				Expect(audit[1].CreatedBy).To(Equal(fun.CreatedByAman))
				Expect(audit[1].CreatedAt).ToNot(BeZero())
			})
		})
	})
})
