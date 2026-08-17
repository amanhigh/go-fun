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

type personGetWatermillStub struct{}

func (personGetWatermillStub) Start(context.Context)    {}
func (personGetWatermillStub) Shutdown(context.Context) {}

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
			Tracer:            tracer,
			PersonHandler:     personHandler,
			EnrollmentHandler: personGetEnrollmentHandlerStub{},
			AdminHandler:      personGetAdminHandlerStub{},
			Watermill:         personGetWatermillStub{},
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
				var err error
				existingPerson, err = personManager.CreatePerson(ctx, fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"})
				Expect(err).ToNot(HaveOccurred())
				Expect(db.Migrator().DropTable(&fun.Person{})).To(Succeed())
				req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/person/"+existingPerson.Id, nil)
				responseRecorder = httptest.NewRecorder()
				router.ServeHTTP(responseRecorder, req)
			})

			It("returns HTTP 500 for an unavailable person table", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
				response := decodePersonGetError(responseRecorder)
				Expect(response["status"]).To(Equal("error"))
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

		Context("Field Validations", func() {
			Context("Person ID Field", func() {
				Context("Allowed Values", func() {
					BeforeEach(func() {
						var err error
						existingPerson, err = personManager.CreatePerson(ctx, fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"})
						Expect(err).ToNot(HaveOccurred())
						req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/person/"+existingPerson.Id+"/audit", nil)
						responseRecorder = httptest.NewRecorder()
						router.ServeHTTP(responseRecorder, req)
					})

					It("accepts a created person ID and returns its audit list", func() {
						Expect(responseRecorder.Code).To(Equal(http.StatusOK))
						Expect(decodePersonAuditResponse(responseRecorder)).To(HaveLen(1))
					})
				})
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
		listPersonsWithResponse := func(query string) (fun.PersonList, *httptest.ResponseRecorder) {
			response := listPersonsResponseOnly(query)
			return decodePersonListResponse(response), response
		}

		Context("Happy Path", func() {
			Context("with no optional query fields", func() {
				var response fun.PersonList

				BeforeEach(func() {
					requests := make([]fun.PersonRequest, 22)
					for i := range requests {
						requests[i] = fun.PersonRequest{Name: "Person " + string(rune('A'+i)), Age: i + 1, Gender: "FEMALE"}
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
				})
			})

			Context("with explicit first and second pages", func() {
				var (
					firstPage, secondPage                 fun.PersonList
					firstPageResponse, secondPageResponse *httptest.ResponseRecorder
				)

				BeforeEach(func() {
					createPersons(
						fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
						fun.PersonRequest{Name: "Bob", Age: 40, Gender: "MALE"},
						fun.PersonRequest{Name: "Carol", Age: 28, Gender: "FEMALE"},
						fun.PersonRequest{Name: "Dave", Age: 50, Gender: "MALE"},
					)
					firstPage, firstPageResponse = listPersonsWithResponse("?sort_by=name&sort-order=asc&offset=0&limit=2")
					secondPage, secondPageResponse = listPersonsWithResponse("?sort_by=name&sort-order=asc&offset=2&limit=2")
				})

				It("returns exact pagination metadata and different records on page two", func() {
					Expect(firstPageResponse.Code).To(Equal(http.StatusOK))
					Expect(secondPageResponse.Code).To(Equal(http.StatusOK))
					Expect(firstPage.Metadata.Offset).To(Equal(0))
					Expect(firstPage.Metadata.Limit).To(Equal(2))
					Expect(firstPage.Metadata.Total).To(Equal(int64(4)))
					Expect(firstPage.Records[0].Name).To(Equal("Ada"))
					Expect(firstPage.Records[1].Name).To(Equal("Bob"))
					Expect(secondPage.Metadata.Offset).To(Equal(2))
					Expect(secondPage.Metadata.Limit).To(Equal(2))
					Expect(secondPage.Records[0].Name).To(Equal("Carol"))
					Expect(secondPage.Records[1].Name).To(Equal("Dave"))
				})
			})

			Context("with a partial name filter", func() {
				var response fun.PersonList

				BeforeEach(func() {
					createPersons(
						fun.PersonRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"},
						fun.PersonRequest{Name: "Ada Byron", Age: 28, Gender: "FEMALE"},
						fun.PersonRequest{Name: "Grace Hopper", Age: 85, Gender: "FEMALE"},
					)
					response = listPersons("?name=Ada")
				})

				It("returns matching records and total metadata", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusOK))
					Expect(response.Records).To(HaveLen(2))
					Expect(response.Records[0].Name).To(Equal("Ada Lovelace"))
					Expect(response.Records[1].Name).To(Equal("Ada Byron"))
					Expect(response.Metadata.Total).To(Equal(int64(2)))
				})
			})

			Context("with an exact gender filter", func() {
				var response fun.PersonList

				BeforeEach(func() {
					createPersons(
						fun.PersonRequest{Name: "Ada", Age: 36, Gender: "FEMALE"},
						fun.PersonRequest{Name: "Grace", Age: 85, Gender: "FEMALE"},
						fun.PersonRequest{Name: "Alan", Age: 42, Gender: "MALE"},
					)
					response = listPersons("?gender=MALE")
				})

				It("returns only the requested gender and total metadata", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusOK))
					Expect(response.Records).To(HaveLen(1))
					Expect(response.Records[0].Name).To(Equal("Alan"))
					Expect(response.Records[0].Gender).To(Equal("MALE"))
					Expect(response.Metadata.Total).To(Equal(int64(1)))
				})
			})

			Context("with combined name and gender filters", func() {
				Context("when both filters match", func() {
					var response fun.PersonList

					BeforeEach(func() {
						createPersons(
							fun.PersonRequest{Name: "Ada Lovelace", Age: 36, Gender: "FEMALE"},
							fun.PersonRequest{Name: "Ada Byron", Age: 28, Gender: "MALE"},
						)
						response = listPersons("?name=Ada&gender=FEMALE")
					})

					It("uses AND semantics", func() {
						Expect(responseRecorder.Code).To(Equal(http.StatusOK))
						Expect(response.Records).To(HaveLen(1))
						Expect(response.Records[0].Name).To(Equal("Ada Lovelace"))
						Expect(response.Metadata.Total).To(Equal(int64(1)))
					})
				})

				Context("when no record satisfies both filters", func() {
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

			Context("with name sorting", func() {
				var (
					ascending, descending                 fun.PersonList
					ascendingResponse, descendingResponse *httptest.ResponseRecorder
				)

				BeforeEach(func() {
					createPersons(
						fun.PersonRequest{Name: "Charlie", Age: 30, Gender: "MALE"},
						fun.PersonRequest{Name: "Alice", Age: 25, Gender: "FEMALE"},
						fun.PersonRequest{Name: "Bob", Age: 40, Gender: "MALE"},
					)
					ascending, ascendingResponse = listPersonsWithResponse("?sort_by=name&sort-order=asc")
					descending, descendingResponse = listPersonsWithResponse("?sort_by=name&sort-order=desc")
				})

				It("sorts names in ascending and descending order", func() {
					Expect(ascendingResponse.Code).To(Equal(http.StatusOK))
					Expect(descendingResponse.Code).To(Equal(http.StatusOK))
					Expect(ascending.Records[0].Name).To(Equal("Alice"))
					Expect(ascending.Records[1].Name).To(Equal("Bob"))
					Expect(ascending.Records[2].Name).To(Equal("Charlie"))
					Expect(descending.Records[0].Name).To(Equal("Charlie"))
					Expect(descending.Records[1].Name).To(Equal("Bob"))
					Expect(descending.Records[2].Name).To(Equal("Alice"))
				})
			})

			Context("with age ascending sorting", func() {
				var response fun.PersonList

				BeforeEach(func() {
					createPersons(
						fun.PersonRequest{Name: "Older", Age: 42, Gender: "MALE"},
						fun.PersonRequest{Name: "Youngest", Age: 18, Gender: "FEMALE"},
						fun.PersonRequest{Name: "Middle", Age: 30, Gender: "FEMALE"},
					)
					response = listPersons("?sort_by=age&sort-order=asc")
				})

				It("returns records in ascending age order", func() {
					Expect(responseRecorder.Code).To(Equal(http.StatusOK))
					Expect(response.Records[0].Name).To(Equal("Youngest"))
					Expect(response.Records[1].Name).To(Equal("Middle"))
					Expect(response.Records[2].Name).To(Equal("Older"))
				})
			})

			Context("with gender sorting", func() {
				var (
					ascending, descending                 fun.PersonList
					ascendingResponse, descendingResponse *httptest.ResponseRecorder
				)

				BeforeEach(func() {
					createPersons(
						fun.PersonRequest{Name: "Male One", Age: 30, Gender: "MALE"},
						fun.PersonRequest{Name: "Female One", Age: 25, Gender: "FEMALE"},
						fun.PersonRequest{Name: "Male Two", Age: 40, Gender: "MALE"},
					)
					ascending, ascendingResponse = listPersonsWithResponse("?sort_by=gender&sort-order=asc")
					descending, descendingResponse = listPersonsWithResponse("?sort_by=gender&sort-order=desc")
				})

				It("sorts genders in ascending and descending order", func() {
					Expect(ascendingResponse.Code).To(Equal(http.StatusOK))
					Expect(descendingResponse.Code).To(Equal(http.StatusOK))
					Expect(ascending.Records[0].Gender).To(Equal("FEMALE"))
					Expect(ascending.Records[1].Gender).To(Equal("MALE"))
					Expect(ascending.Records[2].Gender).To(Equal("MALE"))
					Expect(descending.Records[0].Gender).To(Equal("MALE"))
					Expect(descending.Records[1].Gender).To(Equal("MALE"))
					Expect(descending.Records[2].Gender).To(Equal("FEMALE"))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Offset Field", func() {
				var response fun.PersonList
				Context("Allowed Values", func() {
					Context("with offset 0", func() {
						BeforeEach(func() {
							response = listPersons("?offset=0")
						})

						It("accepts offset 0 and returns a list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(20))
						})
					})

					Context("with a positive offset", func() {
						BeforeEach(func() {
							response = listPersons("?offset=1")
						})

						It("accepts a positive offset and returns a list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(1))
							Expect(response.Metadata.Limit).To(Equal(20))
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
				var response fun.PersonList
				Context("Allowed Values", func() {
					Context("with limit 1", func() {
						BeforeEach(func() { response = listPersons("?limit=1") })

						It("accepts limit 1 and returns a list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(1))
						})
					})

					Context("with limit 100", func() {
						BeforeEach(func() { response = listPersons("?limit=100") })

						It("accepts limit 100 and returns a list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(100))
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
				var response fun.PersonList
				Context("Allowed Values", func() {
					Context("with a partial name", func() {
						BeforeEach(func() { response = listPersons("?name=Ada") })

						It("accepts a partial name and returns a list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(20))
						})
					})

					Context("with a 25-character name", func() {
						BeforeEach(func() {
							response = listPersons("?name=" + strings.Repeat("A", 25))
						})

						It("accepts a 25-character name and returns a list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(20))
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
				var response fun.PersonList
				Context("Allowed Values", func() {
					Context("with MALE", func() {
						BeforeEach(func() { response = listPersons("?gender=MALE") })

						It("accepts MALE and returns a list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(20))
						})
					})

					Context("with FEMALE", func() {
						BeforeEach(func() { response = listPersons("?gender=FEMALE") })

						It("accepts FEMALE and returns a list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(20))
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
				var response fun.PersonList
				Context("Allowed Values", func() {
					Context("with name", func() {
						BeforeEach(func() { response = listPersons("?sort_by=name") })

						It("accepts name and returns a decoded list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(20))
						})
					})

					Context("with age", func() {
						BeforeEach(func() { response = listPersons("?sort_by=age") })

						It("accepts age and returns a decoded list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(20))
						})
					})

					Context("with gender", func() {
						BeforeEach(func() { response = listPersons("?sort_by=gender") })

						It("accepts gender and returns a decoded list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(20))
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
				var response fun.PersonList
				Context("Allowed Values", func() {
					Context("with asc", func() {
						BeforeEach(func() { response = listPersons("?sort-order=asc") })

						It("accepts asc and returns a decoded list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(20))
						})
					})

					Context("with desc", func() {
						BeforeEach(func() { response = listPersons("?sort-order=desc") })

						It("accepts desc and returns a decoded list", func() {
							Expect(responseRecorder.Code).To(Equal(http.StatusOK))
							Expect(response.Metadata.Offset).To(Equal(0))
							Expect(response.Metadata.Limit).To(Equal(20))
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

			It("returns HTTP 500 with the database error as a JSON string", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
				var response string
				Expect(json.Unmarshal(responseRecorder.Body.Bytes(), &response)).To(Succeed())
				Expect(response).ToNot(BeEmpty())
			})
		})
	})
})
