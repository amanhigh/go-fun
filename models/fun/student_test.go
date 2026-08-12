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

var _ = Describe("Student", func() {
	Context("StudentRequest", func() {
		It("should have correct struct fields and tags", func() {
			request := fun.StudentRequest{
				Name:   "John Doe",
				Age:    30,
				Gender: "MALE",
			}

			Expect(request.Name).To(Equal("John Doe"))
			Expect(request.Age).To(Equal(30))
			Expect(request.Gender).To(Equal("MALE"))
		})

		It("should work with female gender", func() {
			request := fun.StudentRequest{
				Name:   "Jane Doe",
				Age:    25,
				Gender: "FEMALE",
			}

			Expect(request.Gender).To(Equal("FEMALE"))
		})
	})

	Context("StudentPath", func() {
		It("should have Id field for URI binding", func() {
			path := fun.StudentPath{
				Id: "abc123",
			}

			Expect(path.Id).To(Equal("abc123"))
		})
	})

	Context("StudentQuery", func() {
		It("should embed Pagination and Sort", func() {
			query := fun.StudentQuery{
				Pagination: common.Pagination{
					Offset: 10,
					Limit:  5,
				},
				Sort: common.Sort{
					SortOrder: common.SortOrderDesc,
				},
				SortBy: "name",
				Name:   "John",
				Gender: "MALE",
			}

			Expect(query.Offset).To(Equal(10))
			Expect(query.Limit).To(Equal(5))
			Expect(query.SortBy).To(Equal("name"))
			Expect(query.SortOrder).To(Equal(common.SortOrderDesc))
			Expect(query.Name).To(Equal("John"))
			Expect(query.Gender).To(Equal("MALE"))
		})
	})

	Context("StudentList", func() {
		It("should contain records and metadata", func() {
			students := []fun.Student{
				{StudentRequest: fun.StudentRequest{Name: "John", Age: 30, Gender: "MALE"}, Id: "1"},
				{StudentRequest: fun.StudentRequest{Name: "Jane", Age: 25, Gender: "FEMALE"}, Id: "2"},
			}

			studentList := fun.StudentList{
				Records: students,
				Metadata: common.PaginatedResponse{
					Total:  2,
					Offset: 0,
					Limit:  20,
				},
			}

			Expect(studentList.Records).To(HaveLen(2))
			Expect(studentList.Records[0].Name).To(Equal("John"))
			Expect(studentList.Records[1].Name).To(Equal("Jane"))
			Expect(studentList.Metadata.Total).To(Equal(int64(2)))
			Expect(studentList.Metadata.Offset).To(Equal(0))
			Expect(studentList.Metadata.Limit).To(Equal(20))
		})
	})

	Context("Student", func() {
		It("should embed StudentRequest and have Id field", func() {
			student := fun.Student{
				StudentRequest: fun.StudentRequest{
					Name:   "John Doe",
					Age:    30,
					Gender: "MALE",
				},
				Id: "abc123",
			}

			Expect(student.Name).To(Equal("John Doe"))
			Expect(student.Age).To(Equal(30))
			Expect(student.Gender).To(Equal("MALE"))
			Expect(student.Id).To(Equal("abc123"))
		})

		Context("BeforeCreate", func() {
			It("should generate 8-character UUID for Id", func() {
				student := &fun.Student{
					StudentRequest: fun.StudentRequest{
						Name:   "Test Student",
						Age:    25,
						Gender: "MALE",
					},
				}

				err := student.BeforeCreate(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(student.Id).NotTo(BeEmpty())
				Expect(student.Id).To(HaveLen(8))
			})

			It("should generate different Ids for different students", func() {
				student1 := &fun.Student{}
				student2 := &fun.Student{}

				err1 := student1.BeforeCreate(nil)
				err2 := student2.BeforeCreate(nil)

				Expect(err1).NotTo(HaveOccurred())
				Expect(err2).NotTo(HaveOccurred())
				Expect(student1.Id).NotTo(Equal(student2.Id))
			})
		})
	})

	Context("CreateStudentAudit", func() {
		It("should create audit from student", func() {
			student := fun.Student{
				StudentRequest: fun.StudentRequest{
					Name:   "John Doe",
					Age:    30,
					Gender: "MALE",
				},
				Id: "abc123",
			}

			audit := fun.CreateStudentAudit(student)

			Expect(audit.Id).To(Equal(student.Id))
			Expect(audit.Name).To(Equal(student.Name))
			Expect(audit.Age).To(Equal(student.Age))
			Expect(audit.Gender).To(Equal(student.Gender))
			// Audit-specific fields should be empty as they're set elsewhere
			Expect(audit.Operation).To(BeEmpty())
			Expect(audit.CreatedBy).To(BeEmpty())
		})
	})

	Context("StudentAudit", func() {
		It("should have all required fields", func() {
			audit := fun.StudentAudit{
				Id:        "abc123",
				Name:      "John Doe",
				Age:       30,
				Gender:    "MALE",
				AuditID:   1,
				Operation: "CREATE",
				CreatedBy: fun.CreatedByAman,
				CreatedAt: time.Now(),
			}

			Expect(audit.Id).To(Equal("abc123"))
			Expect(audit.Name).To(Equal("John Doe"))
			Expect(audit.Age).To(Equal(30))
			Expect(audit.Gender).To(Equal("MALE"))
			Expect(audit.AuditID).To(Equal(uint(1)))
			Expect(audit.Operation).To(Equal("CREATE"))
			Expect(audit.CreatedBy).To(Equal(fun.CreatedByAman))
			Expect(audit.CreatedAt).NotTo(BeZero())
		})

		It("should support different operations", func() {
			operations := []string{"CREATE", "UPDATE", "DELETE"}

			for _, op := range operations {
				audit := fun.StudentAudit{
					Operation: op,
				}
				Expect(audit.Operation).To(Equal(op))
			}
		})
	})

	Context("GORM Hooks Integration", func() {
		Context("Audit Creation Logic", func() {
			It("should create proper audit for CREATE operation", func() {
				student := fun.Student{
					StudentRequest: fun.StudentRequest{
						Name:   "Test User",
						Age:    25,
						Gender: "FEMALE",
					},
					Id: "test123",
				}

				audit := fun.CreateStudentAudit(student)
				audit.Operation = "CREATE"
				audit.CreatedBy = fun.CreatedByAman
				audit.CreatedAt = time.Now()

				Expect(audit.Id).To(Equal("test123"))
				Expect(audit.Name).To(Equal("Test User"))
				Expect(audit.Age).To(Equal(25))
				Expect(audit.Gender).To(Equal("FEMALE"))
				Expect(audit.Operation).To(Equal("CREATE"))
				Expect(audit.CreatedBy).To(Equal(fun.CreatedByAman))
				Expect(audit.CreatedAt).NotTo(BeZero())
			})

			It("should create proper audit for UPDATE operation", func() {
				student := fun.Student{
					StudentRequest: fun.StudentRequest{
						Name:   "Updated User",
						Age:    30,
						Gender: "MALE",
					},
					Id: "update123",
				}

				audit := fun.CreateStudentAudit(student)
				audit.Operation = "UPDATE"
				audit.CreatedBy = fun.CreatedByAman
				audit.CreatedAt = time.Now()

				Expect(audit.Operation).To(Equal("UPDATE"))
			})

			It("should create proper audit for DELETE operation", func() {
				student := fun.Student{
					StudentRequest: fun.StudentRequest{
						Name:   "Deleted User",
						Age:    35,
						Gender: "MALE",
					},
					Id: "delete123",
				}

				audit := fun.CreateStudentAudit(student)
				audit.Operation = "DELETE"
				audit.CreatedBy = fun.CreatedByAman
				audit.CreatedAt = time.Now()

				Expect(audit.Operation).To(Equal("DELETE"))
			})
		})
	})

	Context("Gin Binding Validation", func() {
		var testStudentJSON func(studentJSON string, expectedStatus int)
		var testStudentStruct func(student fun.StudentRequest, expectedStatus int)

		BeforeEach(func() {
			gin.SetMode(gin.TestMode)

			// Register the custom name validator
			if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
				err := v.RegisterValidation("name", funAppCommon.NameValidator)
				Expect(err).NotTo(HaveOccurred())
			}

			testStudentJSON = func(studentJSON string, expectedStatus int) {
				router := gin.New()
				w := httptest.NewRecorder()

				router.POST("/test", func(c *gin.Context) {
					var request fun.StudentRequest
					if err := c.ShouldBindJSON(&request); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, request)
				})

				req, _ := http.NewRequest("POST", "/test", strings.NewReader(studentJSON))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(expectedStatus))
			}

			testStudentStruct = func(student fun.StudentRequest, expectedStatus int) {
				router := gin.New()
				w := httptest.NewRecorder()

				router.POST("/test", func(c *gin.Context) {
					var request fun.StudentRequest
					if err := c.ShouldBindJSON(&request); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, request)
				})

				jsonData, err := json.Marshal(student)
				Expect(err).NotTo(HaveOccurred())

				req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(expectedStatus))
			}
		})

		Context("Valid JSON Binding", func() {
			It("should validate valid StudentRequest", func() {
				validStudent := fun.StudentRequest{
					Name:   "John Smith",
					Age:    25,
					Gender: "MALE",
				}
				testStudentStruct(validStudent, http.StatusOK)
			})

			It("should accept valid female student", func() {
				validStudent := fun.StudentRequest{
					Name:   "Jane Doe",
					Age:    30,
					Gender: "FEMALE",
				}
				testStudentStruct(validStudent, http.StatusOK)
			})
		})

		Context("Invalid JSON Binding", func() {
			It("should reject empty name", func() {
				invalidStudent := fun.StudentRequest{
					Name:   "",
					Age:    25,
					Gender: "MALE",
				}
				testStudentStruct(invalidStudent, http.StatusBadRequest)
			})

			It("should reject age below minimum", func() {
				invalidStudent := fun.StudentRequest{
					Name:   "John Smith",
					Age:    0,
					Gender: "MALE",
				}
				testStudentStruct(invalidStudent, http.StatusBadRequest)
			})

			It("should reject age above maximum", func() {
				invalidStudent := fun.StudentRequest{
					Name:   "John Smith",
					Age:    151,
					Gender: "MALE",
				}
				testStudentStruct(invalidStudent, http.StatusBadRequest)
			})

			It("should reject invalid gender", func() {
				testStudentJSON(`{"name":"John Smith","age":25,"gender":"INVALID"}`, http.StatusBadRequest)
			})

			It("should reject name longer than 25 characters", func() {
				invalidStudent := fun.StudentRequest{
					Name:   "ABCDEFGHIJKLMNOPQRSTUVWXYZ", // 26 characters
					Age:    25,
					Gender: "MALE",
				}
				testStudentStruct(invalidStudent, http.StatusBadRequest)
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
					validStudent := fun.StudentRequest{
						Name:   validName,
						Age:    25,
						Gender: "MALE",
					}
					testStudentStruct(validStudent, http.StatusOK)
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
					invalidStudent := fun.StudentRequest{
						Name:   invalidName,
						Age:    25,
						Gender: "MALE",
					}
					testStudentStruct(invalidStudent, http.StatusBadRequest)
				}
			})
		})

		Context("URI Binding", func() {
			It("should bind StudentPath correctly", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/student/:id", func(c *gin.Context) {
					var path fun.StudentPath
					if err := c.ShouldBindUri(&path); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, path)
				})

				req, _ := http.NewRequest("GET", "/student/abc123", nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				var response fun.StudentPath
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Id).To(Equal("abc123"))
			})
		})

		Context("Query Parameter Binding", func() {
			It("should bind StudentQuery correctly", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/students", func(c *gin.Context) {
					var query fun.StudentQuery
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
				params.Add("sort-order", "asc")
				params.Add("name", "John")
				params.Add("gender", "MALE")

				req, _ := http.NewRequest("GET", "/students?"+params.Encode(), nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				var response fun.StudentQuery
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Offset).To(Equal(10))
				Expect(response.Limit).To(Equal(5))
				Expect(response.SortBy).To(Equal("name"))
				Expect(response.SortOrder).To(Equal(common.SortOrderAsc))
				Expect(response.Name).To(Equal("John"))
				Expect(response.Gender).To(Equal("MALE"))
			})

			It("should reject negative offset", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/students", func(c *gin.Context) {
					var query fun.StudentQuery
					if err := c.ShouldBindQuery(&query); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, query)
				})

				req, _ := http.NewRequest("GET", "/students?offset=-1&limit=5", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should reject limit above 100", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/students", func(c *gin.Context) {
					var query fun.StudentQuery
					if err := c.ShouldBindQuery(&query); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, query)
				})

				req, _ := http.NewRequest("GET", "/students?offset=0&limit=101", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should reject invalid sort_by value", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/students", func(c *gin.Context) {
					var query fun.StudentQuery
					if err := c.ShouldBindQuery(&query); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, query)
				})

				req, _ := http.NewRequest("GET", "/students?sort_by=invalid&order=asc", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should reject invalid gender in query", func() {
				router := gin.New()
				w := httptest.NewRecorder()

				router.GET("/students", func(c *gin.Context) {
					var query fun.StudentQuery
					if err := c.ShouldBindQuery(&query); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, query)
				})

				req, _ := http.NewRequest("GET", "/students?gender=INVALID", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

})
