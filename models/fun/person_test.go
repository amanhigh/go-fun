package fun_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	funAppCommon "github.com/amanhigh/go-fun/components/fun-app/common"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Person", func() {
	Context("PersonRequest", func() {
		It("should have correct struct fields and tags", func() {
			request := fun.PersonRequest{
				Name:   "John Doe",
				Age:    30,
				Gender: "MALE",
			}

			Expect(request.Name).To(Equal("John Doe"))
			Expect(request.Age).To(Equal(30))
			Expect(request.Gender).To(Equal("MALE"))
		})

		It("should work with female gender", func() {
			request := fun.PersonRequest{
				Name:   "Jane Doe",
				Age:    25,
				Gender: "FEMALE",
			}

			Expect(request.Gender).To(Equal("FEMALE"))
		})
	})

	Context("PersonPath", func() {
		It("should have Id field for URI binding", func() {
			path := fun.PersonPath{
				Id: "abc123",
			}

			Expect(path.Id).To(Equal("abc123"))
		})
	})

	Context("PersonQuery", func() {
		It("should embed Pagination and Sort", func() {
			query := fun.PersonQuery{
				Pagination: common.Pagination{
					Offset: 10,
					Limit:  5,
				},
				Sort: common.Sort{
					SortBy: "name",
					Order:  "asc",
				},
				Name:   "John",
				Gender: "MALE",
			}

			Expect(query.Offset).To(Equal(10))
			Expect(query.Limit).To(Equal(5))
			Expect(query.SortBy).To(Equal("name"))
			Expect(query.Order).To(Equal("asc"))
			Expect(query.Name).To(Equal("John"))
			Expect(query.Gender).To(Equal("MALE"))
		})
	})

	Context("PersonList", func() {
		It("should contain records and metadata", func() {
			persons := []fun.Person{
				{PersonRequest: fun.PersonRequest{Name: "John", Age: 30, Gender: "MALE"}, Id: "1"},
				{PersonRequest: fun.PersonRequest{Name: "Jane", Age: 25, Gender: "FEMALE"}, Id: "2"},
			}

			personList := fun.PersonList{
				Records: persons,
				Metadata: common.PaginatedResponse{
					Total: 2,
				},
			}

			Expect(personList.Records).To(HaveLen(2))
			Expect(personList.Records[0].Name).To(Equal("John"))
			Expect(personList.Records[1].Name).To(Equal("Jane"))
			Expect(personList.Metadata.Total).To(Equal(int64(2)))
		})
	})

	Context("Person", func() {
		It("should embed PersonRequest and have Id field", func() {
			person := fun.Person{
				PersonRequest: fun.PersonRequest{
					Name:   "John Doe",
					Age:    30,
					Gender: "MALE",
				},
				Id: "abc123",
			}

			Expect(person.Name).To(Equal("John Doe"))
			Expect(person.Age).To(Equal(30))
			Expect(person.Gender).To(Equal("MALE"))
			Expect(person.Id).To(Equal("abc123"))
		})

		Context("BeforeCreate", func() {
			It("should generate 8-character UUID for Id", func() {
				person := &fun.Person{
					PersonRequest: fun.PersonRequest{
						Name:   "Test Person",
						Age:    25,
						Gender: "MALE",
					},
				}

				err := person.BeforeCreate(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(person.Id).NotTo(BeEmpty())
				Expect(person.Id).To(HaveLen(8))
			})

			It("should generate different Ids for different persons", func() {
				person1 := &fun.Person{}
				person2 := &fun.Person{}

				err1 := person1.BeforeCreate(nil)
				err2 := person2.BeforeCreate(nil)

				Expect(err1).NotTo(HaveOccurred())
				Expect(err2).NotTo(HaveOccurred())
				Expect(person1.Id).NotTo(Equal(person2.Id))
			})
		})
	})

	Context("CreatePersonAudit", func() {
		It("should create audit from person", func() {
			person := fun.Person{
				PersonRequest: fun.PersonRequest{
					Name:   "John Doe",
					Age:    30,
					Gender: "MALE",
				},
				Id: "abc123",
			}

			audit := fun.CreatePersonAudit(person)

			Expect(audit.Id).To(Equal(person.Id))
			Expect(audit.Name).To(Equal(person.Name))
			Expect(audit.Age).To(Equal(person.Age))
			Expect(audit.Gender).To(Equal(person.Gender))
			// Audit-specific fields should be empty as they're set elsewhere
			Expect(audit.Operation).To(BeEmpty())
			Expect(audit.CreatedBy).To(BeEmpty())
		})
	})

	Context("PersonAudit", func() {
		It("should have all required fields", func() {
			audit := fun.PersonAudit{
				Id:        "abc123",
				Name:      "John Doe",
				Age:       30,
				Gender:    "MALE",
				AuditID:   1,
				Operation: "CREATE",
				CreatedBy: "AMAN",
				CreatedAt: time.Now(),
			}

			Expect(audit.Id).To(Equal("abc123"))
			Expect(audit.Name).To(Equal("John Doe"))
			Expect(audit.Age).To(Equal(30))
			Expect(audit.Gender).To(Equal("MALE"))
			Expect(audit.AuditID).To(Equal(uint(1)))
			Expect(audit.Operation).To(Equal("CREATE"))
			Expect(audit.CreatedBy).To(Equal("AMAN"))
			Expect(audit.CreatedAt).NotTo(BeZero())
		})

		It("should support different operations", func() {
			operations := []string{"CREATE", "UPDATE", "DELETE"}

			for _, op := range operations {
				audit := fun.PersonAudit{
					Operation: op,
				}
				Expect(audit.Operation).To(Equal(op))
			}
		})
	})

	Context("GORM Hooks Integration", func() {
		Context("Audit Creation Logic", func() {
			It("should create proper audit for CREATE operation", func() {
				person := fun.Person{
					PersonRequest: fun.PersonRequest{
						Name:   "Test User",
						Age:    25,
						Gender: "FEMALE",
					},
					Id: "test123",
				}

				audit := fun.CreatePersonAudit(person)
				audit.Operation = "CREATE"
				audit.CreatedBy = "AMAN"
				audit.CreatedAt = time.Now()

				Expect(audit.Id).To(Equal("test123"))
				Expect(audit.Name).To(Equal("Test User"))
				Expect(audit.Age).To(Equal(25))
				Expect(audit.Gender).To(Equal("FEMALE"))
				Expect(audit.Operation).To(Equal("CREATE"))
				Expect(audit.CreatedBy).To(Equal("AMAN"))
				Expect(audit.CreatedAt).NotTo(BeZero())
			})

			It("should create proper audit for UPDATE operation", func() {
				person := fun.Person{
					PersonRequest: fun.PersonRequest{
						Name:   "Updated User",
						Age:    30,
						Gender: "MALE",
					},
					Id: "update123",
				}

				audit := fun.CreatePersonAudit(person)
				audit.Operation = "UPDATE"
				audit.CreatedBy = "AMAN"
				audit.CreatedAt = time.Now()

				Expect(audit.Operation).To(Equal("UPDATE"))
			})

			It("should create proper audit for DELETE operation", func() {
				person := fun.Person{
					PersonRequest: fun.PersonRequest{
						Name:   "Deleted User",
						Age:    35,
						Gender: "MALE",
					},
					Id: "delete123",
				}

				audit := fun.CreatePersonAudit(person)
				audit.Operation = "DELETE"
				audit.CreatedBy = "AMAN"
				audit.CreatedAt = time.Now()

				Expect(audit.Operation).To(Equal("DELETE"))
			})
		})
	})

	Context("Gin Binding Validation", func() {
		var testPersonJSON func(personJSON string, expectedStatus int)
		var testPersonStruct func(person fun.PersonRequest, expectedStatus int)

		BeforeEach(func() {
			gin.SetMode(gin.TestMode)

			// Register the custom name validator
			if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
				err := v.RegisterValidation("name", funAppCommon.NameValidator)
				Expect(err).NotTo(HaveOccurred())
			}

			testPersonJSON = func(personJSON string, expectedStatus int) {
				router := gin.New()
				w := httptest.NewRecorder()

				router.POST("/test", func(c *gin.Context) {
					var request fun.PersonRequest
					if err := c.ShouldBindJSON(&request); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, request)
				})

				req, _ := http.NewRequest("POST", "/test", strings.NewReader(personJSON))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(expectedStatus))
			}

			testPersonStruct = func(person fun.PersonRequest, expectedStatus int) {
				router := gin.New()
				w := httptest.NewRecorder()

				router.POST("/test", func(c *gin.Context) {
					var request fun.PersonRequest
					if err := c.ShouldBindJSON(&request); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, request)
				})

				jsonData, err := json.Marshal(person)
				Expect(err).NotTo(HaveOccurred())

				req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(expectedStatus))
			}
		})

		Context("Valid JSON Binding", func() {
			It("should validate valid PersonRequest", func() {
				validPerson := fun.PersonRequest{
					Name:   "John Smith",
					Age:    25,
					Gender: "MALE",
				}
				testPersonStruct(validPerson, http.StatusOK)
			})

			It("should accept valid female person", func() {
				validPerson := fun.PersonRequest{
					Name:   "Jane Doe",
					Age:    30,
					Gender: "FEMALE",
				}
				testPersonStruct(validPerson, http.StatusOK)
			})
		})

		Context("Invalid JSON Binding", func() {
			It("should reject empty name", func() {
				invalidPerson := fun.PersonRequest{
					Name:   "",
					Age:    25,
					Gender: "MALE",
				}
				testPersonStruct(invalidPerson, http.StatusBadRequest)
			})

			It("should reject age below minimum", func() {
				invalidPerson := fun.PersonRequest{
					Name:   "John Smith",
					Age:    0,
					Gender: "MALE",
				}
				testPersonStruct(invalidPerson, http.StatusBadRequest)
			})

			It("should reject age above maximum", func() {
				invalidPerson := fun.PersonRequest{
					Name:   "John Smith",
					Age:    151,
					Gender: "MALE",
				}
				testPersonStruct(invalidPerson, http.StatusBadRequest)
			})

			It("should reject invalid gender", func() {
				testPersonJSON(`{"name":"John Smith","age":25,"gender":"INVALID"}`, http.StatusBadRequest)
			})

			It("should reject name longer than 25 characters", func() {
				invalidPerson := fun.PersonRequest{
					Name:   "ABCDEFGHIJKLMNOPQRSTUVWXYZ", // 26 characters
					Age:    25,
					Gender: "MALE",
				}
				testPersonStruct(invalidPerson, http.StatusBadRequest)
			})
		})

		Context("Custom Name Validator", func() {
			It("should accept valid name characters", func() {
				validNames := []string{
					"John Smith",
					"Mary-Jane Watson",
					"Peter Parker 2",
					"Bruce Wayne",
					"A",
				}

				for _, validName := range validNames {
					validPerson := fun.PersonRequest{
						Name:   validName,
						Age:    25,
						Gender: "MALE",
					}
					testPersonStruct(validPerson, http.StatusOK)
				}
			})

			It("should reject invalid name characters", func() {
				invalidNames := []string{
					"John@Smith",
					"Mary_Jane",
					"Peter#Parker",
					"Bruce$Wayne",
					"Tony*Stark",
				}

				for _, invalidName := range invalidNames {
					invalidPerson := fun.PersonRequest{
						Name:   invalidName,
						Age:    25,
						Gender: "MALE",
					}
					testPersonStruct(invalidPerson, http.StatusBadRequest)
				}
			})
		})

		Context("URI Binding", func() {
			It("should bind PersonPath correctly", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/person/:id", func(c *gin.Context) {
					var path fun.PersonPath
					if err := c.ShouldBindUri(&path); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, path)
				})

				req, _ := http.NewRequest("GET", "/person/abc123", nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				var response fun.PersonPath
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Id).To(Equal("abc123"))
			})
		})

		Context("Query Parameter Binding", func() {
			It("should bind PersonQuery correctly", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/persons", func(c *gin.Context) {
					var query fun.PersonQuery
					if err := c.ShouldBindQuery(&query); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, query)
				})

				params := url.Values{}
				params.Add("offset", "10")
				params.Add("limit", "5")
				params.Add("sort_by", "name")
				params.Add("order", "asc")
				params.Add("name", "John")
				params.Add("gender", "MALE")

				req, _ := http.NewRequest("GET", "/persons?"+params.Encode(), nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				var response fun.PersonQuery
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Offset).To(Equal(10))
				Expect(response.Limit).To(Equal(5))
				Expect(response.SortBy).To(Equal("name"))
				Expect(response.Order).To(Equal("asc"))
				Expect(response.Name).To(Equal("John"))
				Expect(response.Gender).To(Equal("MALE"))
			})

			It("should reject negative offset", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/persons", func(c *gin.Context) {
					var query fun.PersonQuery
					if err := c.ShouldBindQuery(&query); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, query)
				})

				req, _ := http.NewRequest("GET", "/persons?offset=-1&limit=5", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should reject limit above 10", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/persons", func(c *gin.Context) {
					var query fun.PersonQuery
					if err := c.ShouldBindQuery(&query); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, query)
				})

				req, _ := http.NewRequest("GET", "/persons?offset=0&limit=11", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should reject invalid sort_by value", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/persons", func(c *gin.Context) {
					var query fun.PersonQuery
					if err := c.ShouldBindQuery(&query); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, query)
				})

				req, _ := http.NewRequest("GET", "/persons?sort_by=invalid&order=asc", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should reject invalid gender in query", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/persons", func(c *gin.Context) {
					var query fun.PersonQuery
					if err := c.ShouldBindQuery(&query); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, query)
				})

				req, _ := http.NewRequest("GET", "/persons?gender=INVALID", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

})
