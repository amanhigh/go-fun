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
	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/noop"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type personGetEnrollmentHandlerStub struct{}

func (personGetEnrollmentHandlerStub) CreateEnrollment(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
func (personGetEnrollmentHandlerStub) GetEnrollment(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

type personGetAdminHandlerStub struct{}

func (personGetAdminHandlerStub) Stop(c *gin.Context) { c.Status(http.StatusNotImplemented) }

func decodePersonGetResponse(responseRecorder *httptest.ResponseRecorder) fun.Person {
	var response fun.Person
	Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
	return response
}

func decodePersonListResponse(responseRecorder *httptest.ResponseRecorder) fun.PersonList {
	var response fun.PersonList
	Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
	return response
}

func decodePersonAuditResponse(responseRecorder *httptest.ResponseRecorder) []fun.PersonAudit {
	var response []fun.PersonAudit
	Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
	return response
}

func decodePersonGetError(responseRecorder *httptest.ResponseRecorder) map[string]any {
	var response map[string]any
	Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
	return response
}

var _ = Describe("Person Handler Integration - GET Tests", func() {
	var (
		ctx              context.Context
		db               *gorm.DB
		dbSQL            *sql.DB
		personManager    manager.PersonManagerInterface
		router           *gin.Engine
		existingPerson   fun.Person
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

		tracer := otel.Tracer("fun-app-person-get-handler-test")
		personManager = manager.NewPersonManager(dao.NewPersonDao(util.NewBaseDbRepository(db)), tracer)
		meter := noop.NewMeterProvider().Meter("fun-app-person-get-handler-test")
		createCounter, err := meter.Int64Counter("get_test_create_person")
		Expect(err).ToNot(HaveOccurred())
		personCounter, err := meter.Int64UpDownCounter("get_test_person_count")
		Expect(err).ToNot(HaveOccurred())
		personCreateTime, err := meter.Float64Histogram("get_test_person_create_time")
		Expect(err).ToNot(HaveOccurred())
		personHandler := &handlers.PersonHandlerImpl{
			Manager:          personManager,
			Tracer:           tracer,
			CreateCounter:    createCounter,
			PersonCounter:    personCounter,
			PersonCreateTime: personCreateTime,
		}

		lifecycle := &handlers.FunAppServerLifecycle{
			PersonHandler:     personHandler,
			EnrollmentHandler: personGetEnrollmentHandlerStub{},
			AdminHandler:      personGetAdminHandlerStub{},
		}
		router = util.CreateTestGinRouter()
		lifecycle.RegisterRoutes(router)
	})

	AfterEach(func() {
		if dbSQL != nil {
			Expect(dbSQL.Close()).To(Succeed())
		}
	})

	Describe("GET /v1/person/:id", func() {
		Context("Happy Path", func() {
			BeforeEach(func() {
				var err error
				existingPerson, err = personManager.CreatePerson(ctx, fun.PersonRequest{
					Name: "Ada Lovelace", Age: 36, Gender: "FEMALE",
				})
				Expect(err).ToNot(HaveOccurred())
				req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/person/"+existingPerson.Id, nil)
				responseRecorder = httptest.NewRecorder()
				router.ServeHTTP(responseRecorder, req)
			})

			It("returns all persisted person fields", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
				response := decodePersonGetResponse(responseRecorder)
				Expect(response.Id).To(Equal(existingPerson.Id))
				Expect(response.Name).To(Equal(existingPerson.Name))
				Expect(response.Age).To(Equal(existingPerson.Age))
				Expect(response.Gender).To(Equal(existingPerson.Gender))
			})
		})

		Context("Field Validations", func() {
			Context("Person ID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						var err error
						existingPerson, err = personManager.CreatePerson(ctx, fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"})
						Expect(err).ToNot(HaveOccurred())
						req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/person/"+existingPerson.Id, nil)
						responseRecorder = httptest.NewRecorder()
						router.ServeHTTP(responseRecorder, req)
					})

					It("accepts a created and persisted person ID", func() {
						Expect(responseRecorder.Code).To(Equal(http.StatusOK))
						response := decodePersonGetResponse(responseRecorder)
						Expect(response.Id).To(Equal(existingPerson.Id))
					})
				})

				Context("Bad Values", func() {
					BeforeEach(func() {
						req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/person/missing-person", nil)
						responseRecorder = httptest.NewRecorder()
						router.ServeHTTP(responseRecorder, req)
					})

					It("returns a meaningful JSend not-found failure", func() {
						Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
						response := decodePersonGetError(responseRecorder)
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
				Expect(db.Migrator().DropTable(&fun.Person{})).To(Succeed())
				req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/person/fixed-id", nil)
				responseRecorder = httptest.NewRecorder()
				router.ServeHTTP(responseRecorder, req)
			})

			It("returns HTTP 500 for an unavailable person table", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
				response := decodePersonGetError(responseRecorder)
				Expect(response["status"]).To(Equal("error"))
				Expect(response["message"]).ToNot(BeEmpty())
				Expect(response["code"]).To(Equal(float64(http.StatusInternalServerError)))
			})
		})
	})

	Describe("GET /v1/person/:id/audit", func() {
		Context("Happy Path", func() {
			var audit []fun.PersonAudit

			BeforeEach(func() {
				var err error
				existingPerson, err = personManager.CreatePerson(ctx, fun.PersonRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"})
				Expect(err).ToNot(HaveOccurred())
				Expect(personManager.UpdatePerson(ctx, existingPerson.Id, fun.PersonRequest{Name: "Grace Hopper", Age: 85, Gender: "FEMALE"})).ToNot(HaveOccurred())
				req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/person/"+existingPerson.Id+"/audit", nil)
				responseRecorder = httptest.NewRecorder()
				router.ServeHTTP(responseRecorder, req)
				audit = decodePersonAuditResponse(responseRecorder)
			})

			It("returns ordered CREATE and UPDATE audit records", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
				Expect(audit).To(HaveLen(2))
				Expect(audit[0].AuditID).ToNot(BeZero())
				Expect(audit[0].Id).To(Equal(existingPerson.Id))
				Expect(audit[0].Name).To(Equal("Ada Lovelace"))
				Expect(audit[0].Age).To(Equal(36))
				Expect(audit[0].Gender).To(Equal("FEMALE"))
				Expect(audit[0].Operation).To(Equal("CREATE"))
				Expect(audit[0].CreatedBy).To(Equal(fun.CreatedByAman))
				Expect(audit[0].CreatedAt).ToNot(BeZero())
				Expect(audit[1].AuditID).ToNot(BeZero())
				Expect(audit[1].Id).To(Equal(existingPerson.Id))
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
				Expect(db.Migrator().DropTable(&fun.PersonAudit{})).To(Succeed())
				req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/person/missing-person/audit", nil)
				responseRecorder = httptest.NewRecorder()
				router.ServeHTTP(responseRecorder, req)
			})

			It("returns HTTP 500 for an unavailable audit table", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
				response := decodePersonGetError(responseRecorder)
				Expect(response["status"]).To(Equal("error"))
				Expect(response["message"]).ToNot(BeEmpty())
				Expect(response["code"]).To(Equal(float64(http.StatusInternalServerError)))
			})
		})
	})

	Describe("GET /v1/person/", func() {
		createPersons := func(requests ...fun.PersonRequest) {
			for _, request := range requests {
				_, err := personManager.CreatePerson(ctx, request)
				Expect(err).ToNot(HaveOccurred())
			}
		}
		listPersonsResponseOnly := func(query string) *httptest.ResponseRecorder {
			req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/person/"+query, nil)
			responseRecorder = httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			return responseRecorder
		}
		listPersons := func(query string) fun.PersonList {
			return decodePersonListResponse(listPersonsResponseOnly(query))
		}

		Context("Happy Path", func() {
			Context("with no optional query fields", func() {
				var response fun.PersonList

				BeforeEach(func() {
					requests := make([]fun.PersonRequest, 22)
					for i := range requests {
						// Descending name order: Person V, Person U, ..., Person A
						requests[i] = fun.PersonRequest{Name: "Person " + string(rune('V'-i)), Age: 22 - i, Gender: "FEMALE"}
					}
					createPersons(requests...)
					response = listPersons("")
				})

				It("returns a raw list with default pagination and total metadata", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusOK))
					Expect(response.Records).To(HaveLen(20))
					Expect(response.Metadata.Offset).To(Equal(0))
					Expect(response.Metadata.Limit).To(Equal(20))
					Expect(response.Metadata.Total).To(Equal(int64(22)))
					Expect(response.Records[0].Name).To(Equal("Person A"))
					Expect(response.Records[19].Name).To(Equal("Person T"))
				})
			})

			Context("with combined name and gender filters proving AND semantics with empty result", func() {
				var response fun.PersonList

				BeforeEach(func() {
					createPersons(fun.PersonRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"})
					response = listPersons("?name=Ada&gender=MALE")
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
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Bob", Age: 40, Gender: "MALE"},
								fun.PersonRequest{Name: "Carol", Age: 28, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Dave", Age: 50, Gender: "MALE"},
							)
							response = listPersons("?offset=0&limit=2&sort_by=name&sort-order=asc")
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
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Bob", Age: 40, Gender: "MALE"},
								fun.PersonRequest{Name: "Carol", Age: 28, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Dave", Age: 50, Gender: "MALE"},
							)
							response = listPersons("?offset=2&limit=2&sort_by=name&sort-order=asc")
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
						responseRecorder = listPersonsResponseOnly("?offset=-1")
					})

					It("returns a minimum Offset validation error", func() {
						util.AssertError(responseRecorder, "Offset", "min")
					})
				})
			})

			Context("Limit Field", func() {
				Context("Allowed Values", func() {
					Context("with limit 1", func() {
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Bob", Age: 40, Gender: "MALE"},
								fun.PersonRequest{Name: "Carol", Age: 28, Gender: "FEMALE"},
							)
							response = listPersons("?limit=1&sort_by=name&sort-order=asc")
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
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Bob", Age: 40, Gender: "MALE"},
								fun.PersonRequest{Name: "Carol", Age: 28, Gender: "FEMALE"},
							)
							response = listPersons("?limit=100&sort_by=name&sort-order=asc")
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
						BeforeEach(func() { responseRecorder = listPersonsResponseOnly("?limit=0") })

						It("returns a minimum Limit validation error", func() {
							util.AssertError(responseRecorder, "Limit", "min")
						})
					})

					Context("with limit 101", func() {
						BeforeEach(func() { responseRecorder = listPersonsResponseOnly("?limit=101") })

						It("returns a maximum Limit validation error", func() {
							util.AssertError(responseRecorder, "Limit", "max")
						})
					})
				})
			})

			Context("Name Field", func() {
				Context("Allowed Values", func() {
					Context("with a partial name", func() {
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Ada Byron", Age: 28, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Grace Hopper", Age: 85, Gender: "FEMALE"},
							)
							response = listPersons("?name=Ada&sort_by=name")
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
						var response fun.PersonList

						BeforeEach(func() {
							name25 := strings.Repeat("A", 25)
							createPersons(
								fun.PersonRequest{Name: name25, Age: 30, Gender: "MALE"},
								fun.PersonRequest{Name: "Bob", Age: 40, Gender: "MALE"},
							)
							response = listPersons("?name=" + name25)
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
						BeforeEach(func() { responseRecorder = listPersonsResponseOnly("?name=A%2AB") })

						It("returns a Name character validation error", func() {
							util.AssertError(responseRecorder, "Name", "name")
						})
					})

					Context("with a 26-character name", func() {
						BeforeEach(func() {
							responseRecorder = listPersonsResponseOnly("?name=" + strings.Repeat("A", 26))
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
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Alan", Age: 42, Gender: "MALE"},
								fun.PersonRequest{Name: "Bob", Age: 30, Gender: "MALE"},
								fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
							)
							response = listPersons("?gender=MALE")
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
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Grace", Age: 85, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Bob", Age: 30, Gender: "MALE"},
							)
							response = listPersons("?gender=FEMALE")
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
						BeforeEach(func() { responseRecorder = listPersonsResponseOnly("?gender=OTHER") })

						It("returns an equality Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "eq")
						})
					})

					Context("with lowercase female", func() {
						BeforeEach(func() { responseRecorder = listPersonsResponseOnly("?gender=female") })

						It("returns an equality Gender validation error", func() {
							util.AssertError(responseRecorder, "Gender", "eq")
						})
					})
				})
			})

			Context("SortBy Field", func() {
				Context("Allowed Values", func() {
					Context("with name", func() {
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Charlie", Age: 30, Gender: "MALE"},
								fun.PersonRequest{Name: "Alice", Age: 25, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Bob", Age: 40, Gender: "MALE"},
							)
							response = listPersons("?sort_by=name")
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
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Older", Age: 42, Gender: "MALE"},
								fun.PersonRequest{Name: "Youngest", Age: 18, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Middle", Age: 30, Gender: "FEMALE"},
							)
							response = listPersons("?sort_by=age")
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
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Male One", Age: 30, Gender: "MALE"},
								fun.PersonRequest{Name: "Female One", Age: 25, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Male Two", Age: 40, Gender: "MALE"},
							)
							response = listPersons("?sort_by=gender")
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
					BeforeEach(func() { responseRecorder = listPersonsResponseOnly("?sort_by=invalid") })

					It("returns an equality SortBy validation error", func() {
						util.AssertError(responseRecorder, "SortBy", "eq")
					})
				})
			})

			Context("SortOrder Field", func() {
				Context("Allowed Values", func() {
					Context("with asc", func() {
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Charlie", Age: 30, Gender: "MALE"},
								fun.PersonRequest{Name: "Alice", Age: 25, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Bob", Age: 40, Gender: "MALE"},
							)
							response = listPersons("?sort_by=name&sort-order=asc")
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
						var response fun.PersonList

						BeforeEach(func() {
							createPersons(
								fun.PersonRequest{Name: "Charlie", Age: 30, Gender: "MALE"},
								fun.PersonRequest{Name: "Alice", Age: 25, Gender: "FEMALE"},
								fun.PersonRequest{Name: "Bob", Age: 40, Gender: "MALE"},
							)
							response = listPersons("?sort_by=name&sort-order=desc")
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
					BeforeEach(func() { responseRecorder = listPersonsResponseOnly("?sort-order=invalid") })

					It("returns a oneof SortOrder validation error", func() {
						util.AssertError(responseRecorder, "SortOrder", "oneof")
					})
				})
			})
		})

		Context("Errors", func() {
			BeforeEach(func() {
				Expect(db.Migrator().DropTable(&fun.Person{})).To(Succeed())
				responseRecorder = listPersonsResponseOnly("")
			})

			// Intentionally assert only status code: the raw response body format
			// (raw string vs JSend envelope) is an implementation detail not yet stabilised.
			It("returns HTTP 500 when the database table is missing", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
