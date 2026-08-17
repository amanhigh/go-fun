package it_test

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	. "github.com/amanhigh/go-fun/common/clients"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Entr: http://eradman.com/entrproject/entr.1.html
// -s (use Shell), -c (Clear), Space/Q to Run, Quit.
// Watch Mode: find `git rev-parse --show-toplevel` | entr -s "date +%M:%S;ginkgo $PWD | grep Pending"
var _ = Describe("Student Integration Test", func() {
	// FIXME: Merge with enrollment_test.go and limit integration tests to only important ones
	const (
		invalidNameValue  = "A*B"
		expectedGenderErr = "FEMALE"
		nameFieldErr      = "Name"
		maxValidationTag  = "max"
		nameSortField     = "name"
		reqdValidationTag = "required"
	)
	var (
		request fun.StudentRequest

		name        = "Amanpreet Singh"
		maxName     = strings.Repeat("A", 31)
		age         = 31
		gender      = "MALE"
		err         common.HttpError
		testCtx     = context.Background()
		expectedErr string
	)

	BeforeEach(func() {
		expectedErr = "Bad Request"
		request = fun.StudentRequest{
			Name:   name,
			Age:    age,
			Gender: gender,
		}
	})

	Context("Create", func() {
		var (
			createdStudent fun.Student
			auditUser     = "AMAN"
		)
		BeforeEach(func() {
			createdStudent, err = client.StudentService.CreateStudent(testCtx, request)
			Expect(err).ToNot(HaveOccurred())
			Expect(createdStudent.Id).Should(Not(BeEmpty()))
		})
		AfterEach(func() {
			// Delete Student
			err = client.StudentService.DeleteStudent(testCtx, createdStudent.Id)
			Expect(err).ToNot(HaveOccurred())

			// Delete Audit
			auditList, listErr := client.StudentService.ListStudentAudit(testCtx, createdStudent.Id)
			Expect(listErr).ShouldNot(HaveOccurred())
			Expect(auditList).To(HaveLen(2))
		})

		It("should create & get student", func() {
			student, getErr := client.StudentService.GetStudent(testCtx, createdStudent.Id)
			Expect(getErr).ShouldNot(HaveOccurred())
			Expect(student).Should(Not(BeNil()))

			// Match Student Fields
			Expect(student.Id).To(Equal(createdStudent.Id))
			Expect(student.Name).To(Equal(name))
			Expect(student.Age).To(Equal(age))
			Expect(student.Gender).To(Equal(gender))
		})

		It("should generate Audit", func() {
			// List Audit
			auditList, auditErr := client.StudentService.ListStudentAudit(testCtx, createdStudent.Id)
			Expect(auditErr).ShouldNot(HaveOccurred())

			// Check Audit
			Expect(auditList).To(HaveLen(1))
			audit := auditList[0]
			Expect(audit.Id).To(Equal(createdStudent.Id))
			Expect(audit.Name).To(Equal(name))
			Expect(audit.Age).To(Equal(age))
			Expect(audit.Gender).To(Equal(gender))

			Expect(audit.Operation).To(Equal("CREATE"))
			Expect(audit.CreatedBy).To(Equal(auditUser))
			Expect(audit.CreatedAt).Should(Not(BeNil()))
		})

		Context("Update", func() {
			var (
				updateRequest fun.StudentRequest
				updatedStudent fun.Student
			)
			BeforeEach(func() {
				updateRequest = fun.StudentRequest{
					Name:   "Jenny",
					Age:    25,
					Gender: "FEMALE",
				}
				updatedStudent, err = client.StudentService.CreateStudent(testCtx, request)
				Expect(err).ShouldNot(HaveOccurred())
			})

			AfterEach(func() {
				err = client.StudentService.DeleteStudent(testCtx, updatedStudent.Id)
				Expect(err).ToNot(HaveOccurred())
			})

			Context("Success", func() {
				BeforeEach(func() {
					updateErr := client.StudentService.UpdateStudent(testCtx, updatedStudent.Id, updateRequest)
					Expect(updateErr).ShouldNot(HaveOccurred())
				})

				It("should update student", func() {
					// Fetch Update Student
					student, getErr := client.StudentService.GetStudent(testCtx, updatedStudent.Id)
					Expect(getErr).ShouldNot(HaveOccurred())

					// MatchFields
					Expect(student.Id).To(Equal(updatedStudent.Id))
					Expect(student.Name).To(Equal(updateRequest.Name))
					Expect(student.Age).To(Equal(updateRequest.Age))
					Expect(student.Gender).To(Equal(updateRequest.Gender))
				})

				It("should generate Audit", func() {
					// List Audit
					auditList, auditErr := client.StudentService.ListStudentAudit(testCtx, updatedStudent.Id)
					Expect(auditErr).ShouldNot(HaveOccurred())

					// Check Audit
					Expect(auditList).To(HaveLen(2))
					audit := auditList[1]
					Expect(audit.Id).To(Equal(updatedStudent.Id))
					Expect(audit.Name).To(Equal(updateRequest.Name))
					Expect(audit.Age).To(Equal(updateRequest.Age))
					Expect(audit.Gender).To(Equal(updateRequest.Gender))

					Expect(audit.Operation).To(Equal("UPDATE"))
					Expect(audit.CreatedBy).To(Equal(auditUser))
					Expect(audit.CreatedAt).Should(Not(BeNil()))
				})
			})

			Context("Bad Requests", func() {
				AfterEach(func() {
					err = client.StudentService.UpdateStudent(testCtx, updatedStudent.Id, updateRequest)
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(expectedErr))
					Expect(err.Code()).To(Equal(http.StatusBadRequest))
				})

				It("should fail for missing Name", func() {
					updateRequest.Name = ""
					expectedErr = reqdValidationTag
				})

				It("should fail for invalid Name", func() {
					updateRequest.Name = invalidNameValue
					expectedErr = nameFieldErr
				})

				It("should fail for max Name", func() {
					updateRequest.Name = maxName
					expectedErr = maxValidationTag
				})

				It("should fail for missing Age", func() {
					updateRequest.Age = 0
					expectedErr = "Age"
				})

				It("should fail for invalid Age", func() {
					updateRequest.Age = -1
					expectedErr = "min"
				})

				It("should fail for max Age", func() {
					updateRequest.Age = 200
					expectedErr = maxValidationTag
				})

				It("should fail for missing Gender", func() {
					updateRequest.Gender = ""
					expectedErr = reqdValidationTag
				})

				It("should fail for invalid Gender", func() {
					updateRequest.Gender = "GENDER"
					expectedErr = expectedGenderErr
				})

			})
		})

		Context("Search", func() {
			var (
				offset      = 0
				limit       = 5
				total       = 15
				studentQuery fun.StudentQuery
				names       = []string{"Jane", "Sardar", "Rahul"}
				genders     = []string{"FEMALE", "MALE", "MALE"}
			)

			BeforeEach(func() {
				// Create 15 Students
				for i := range total {
					request.Name = names[i%3] + strconv.Itoa(i)
					request.Gender = genders[i%3]
					_, err = client.StudentService.CreateStudent(testCtx, request)
					Expect(err).ToNot(HaveOccurred())
				}

				// Init Student Query
				studentQuery = fun.StudentQuery{
					Pagination: common.Pagination{
						Offset: offset,
						Limit:  limit,
					},
				}
			})

			AfterEach(func() {
				// Find Record By Names and Delete using UUID
				for i, name := range names {
					studentQuery.Name = name
					studentQuery.Gender = genders[i]
					studentQuery.Limit = 10
					studentQuery.Offset = 0
					studentList, listErr := client.StudentService.ListStudent(testCtx, studentQuery)
					Expect(listErr).ToNot(HaveOccurred())

					// Delete all Records of Name
					for _, student := range studentList.Records {
						err = client.StudentService.DeleteStudent(testCtx, student.Id)
						Expect(err).ToNot(HaveOccurred())
					}
				}
			})

			It("should get all students upto page Limit", func() {
				var studentList fun.StudentList
				studentList, err = client.StudentService.ListStudent(testCtx, studentQuery)
				Expect(err).ToNot(HaveOccurred())

				// Student Count should be same as Page Limit
				Expect(studentList.Records).To(HaveLen(limit))
				Expect(studentList.Metadata.Total).To(BeNumerically(">=", total))
				Expect(studentList.Metadata.Offset).To(Equal(0))
				Expect(studentList.Metadata.Limit).To(Equal(limit))
			})

			It("should fetch second Page", func() {
				var studentList fun.StudentList
				studentQuery.Offset = limit
				studentList, err = client.StudentService.ListStudent(testCtx, studentQuery)

				Expect(err).ToNot(HaveOccurred())
				Expect(studentList.Records).To(HaveLen(limit))
				Expect(studentList.Metadata.Offset).To(Equal(limit))
				Expect(studentList.Metadata.Limit).To(Equal(limit))
			})

			It("should search by Name", func() {
				var studentList fun.StudentList
				studentQuery.Name = names[0]
				studentList, err = client.StudentService.ListStudent(testCtx, studentQuery)

				Expect(err).ToNot(HaveOccurred())
				Expect(studentList.Records).To(HaveLen(limit))
				Expect(studentList.Metadata.Total).To(BeEquivalentTo(5))
				Expect(studentList.Metadata.Offset).To(Equal(0))
				Expect(studentList.Metadata.Limit).To(Equal(limit))
			})

			It("should search by Gender", func() {
				var studentList fun.StudentList
				studentQuery.Gender = genders[1]
				studentList, err = client.StudentService.ListStudent(testCtx, studentQuery)

				Expect(err).ToNot(HaveOccurred())
				Expect(studentList.Records).To(HaveLen(limit))
				Expect(studentList.Metadata.Total).To(BeEquivalentTo(11))
				Expect(studentList.Metadata.Offset).To(Equal(0))
				Expect(studentList.Metadata.Limit).To(Equal(limit))
			})

			It("should search by Name & Gender", func() {
				var studentList fun.StudentList
				studentQuery.Name = names[0]
				studentQuery.Gender = genders[1]
				studentList, err = client.StudentService.ListStudent(testCtx, studentQuery)

				Expect(err).ToNot(HaveOccurred())
				Expect(studentList.Records).To(BeEmpty())
				Expect(studentList.Metadata.Total).To(BeEquivalentTo(0))
				Expect(studentList.Metadata.Offset).To(Equal(0))
				Expect(studentList.Metadata.Limit).To(Equal(limit))
			})

			Context("Sort", func() {

				It("should sort by Name in ascending order", func() {
					var studentList fun.StudentList
					studentQuery.SortBy = nameSortField
					studentQuery.SortOrder = common.SortOrderAsc
					studentList, err = client.StudentService.ListStudent(testCtx, studentQuery)

					Expect(err).ToNot(HaveOccurred())
					Expect(studentList.Records).To(HaveLen(limit))
					Expect(studentList.Metadata.Offset).To(Equal(0))
					Expect(studentList.Metadata.Limit).To(Equal(limit))
					// Check if the records are sorted in ascending order by name
					for i := 0; i < len(studentList.Records)-1; i++ {
						cur := studentList.Records[i].Name
						next := studentList.Records[i+1].Name
						Expect(cur <= next).To(BeTrue())
					}
				})

				It("should sort by Name in descending order", func() {
					var studentList fun.StudentList
					studentQuery.SortBy = nameSortField
					studentQuery.SortOrder = common.SortOrderDesc
					studentList, err = client.StudentService.ListStudent(testCtx, studentQuery)

					Expect(err).ToNot(HaveOccurred())
					Expect(studentList.Records).To(HaveLen(limit))
					Expect(studentList.Metadata.Offset).To(Equal(0))
					Expect(studentList.Metadata.Limit).To(Equal(limit))

					// Check if the records are sorted in descending order by name
					for i := 0; i < len(studentList.Records)-1; i++ {
						cur := studentList.Records[i].Name
						next := studentList.Records[i+1].Name
						Expect(cur >= next).To(BeTrue())
					}
				})

				It("should sort by Gender in ascending order", func() {
					var studentList fun.StudentList
					studentQuery.SortBy = "gender"
					studentQuery.SortOrder = common.SortOrderAsc
					studentList, err = client.StudentService.ListStudent(testCtx, studentQuery)

					Expect(err).ToNot(HaveOccurred())
					Expect(studentList.Records).To(HaveLen(limit))
					Expect(studentList.Metadata.Offset).To(Equal(0))
					Expect(studentList.Metadata.Limit).To(Equal(limit))

					// Check if the records are sorted in ascending order by gender
					for i := 0; i < len(studentList.Records)-1; i++ {
						cur := studentList.Records[i].Gender
						next := studentList.Records[i+1].Gender
						Expect(cur <= next).To(BeTrue())
					}
				})

				It("should sort by Gender in descending order", func() {
					var studentList fun.StudentList
					studentQuery.SortBy = "gender"
					studentQuery.SortOrder = common.SortOrderDesc
					studentList, err = client.StudentService.ListStudent(testCtx, studentQuery)

					Expect(err).ToNot(HaveOccurred())
					Expect(studentList.Records).To(HaveLen(limit))
					Expect(studentList.Metadata.Offset).To(Equal(0))
					Expect(studentList.Metadata.Limit).To(Equal(limit))

					// Check if the records are sorted in descending order by gender
					for i := 0; i < len(studentList.Records)-1; i++ {
						cur := studentList.Records[i].Gender
						next := studentList.Records[i+1].Gender
						Expect(cur >= next).To(BeTrue())
					}
				})
			})

			Context("Bad Requests", func() {
				AfterEach(func() {
					_, err = client.StudentService.ListStudent(testCtx, studentQuery)
					Expect(err).Should(HaveOccurred())
					Expect(err.Code()).To(Equal(http.StatusBadRequest))
					Expect(err.Error()).To(ContainSubstring(expectedErr))

					// Pollutes AfterEach Cleanup so Reset
					studentQuery.SortOrder = common.SortOrderNone
					studentQuery.SortBy = ""
				})

				It("should fail for invalid Offset", func() {
					studentQuery.Offset = -1
					expectedErr = "Offset"
				})

				It("should fail for Lower Limit", func() {
					studentQuery.Limit = 0
					expectedErr = "min (1)"
				})

				It("should fail for Max Limit", func() {
					studentQuery.Limit = 101
					expectedErr = "max (100)"
				})

				It("should fail for invalid Name", func() {
					studentQuery.Name = invalidNameValue
					expectedErr = nameFieldErr
				})

				It("should fail for max Name", func() {
					studentQuery.Name = maxName
					expectedErr = maxValidationTag
				})

				It("should fail for invalid Gender", func() {
					studentQuery.Gender = "OTHER"
					expectedErr = expectedGenderErr
				})

				It("should fail for invalid SortBy", func() {
					studentQuery.SortBy = "invalid"
					expectedErr = "SortBy"
				})

				It("should fail for invalid Order", func() {
					studentQuery.SortBy = nameSortField
					studentQuery.SortOrder = common.SortOrder("invalid")
					expectedErr = "asc"
				})
			})
		})

		Context("Bad Requests", func() {
			AfterEach(func() {
				_, err = client.StudentService.CreateStudent(testCtx, request)

				Expect(err).Should(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusBadRequest))
				Expect(err.Error()).To(ContainSubstring(expectedErr))
			})

			It("should fail for missing Name", func() {
				request.Name = ""
				expectedErr = reqdValidationTag
			})

			It("should fail for invalid Name", func() {
				request.Name = invalidNameValue
				expectedErr = nameFieldErr
			})

			It("should fail for max Name", func() {
				request.Name = maxName
				expectedErr = maxValidationTag
			})

			It("should fail for minimum Age", func() {
				request.Age = 0
				expectedErr = "Age"
			})

			It("should fail for max Age", func() {
				request.Age = 200
				expectedErr = maxValidationTag
			})

			It("should fail for missing Gender", func() {
				request.Gender = ""
				expectedErr = "Gender"
			})

			It("should fail for invalid Gender", func() {
				request.Gender = "OTHER"
				expectedErr = expectedGenderErr
			})
		})
	})

	Context("Bad Requests", func() {
		var (
			emptyId   = ""
			missingId = "missing-id"
		)

		Context("Empty Id", func() {
			AfterEach(func() {
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(common.ErrNotFound))
				Expect(err.Code()).To(Equal(404))
			})

			It("should fail for delete", func() {
				err = client.StudentService.DeleteStudent(testCtx, emptyId)
			})
		})

		Context("Missing Id", func() {
			AfterEach(func() {
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(common.ErrNotFound))
			})

			It("should fail for get", func() {
				_, err = client.StudentService.GetStudent(testCtx, missingId)
			})

			It("should fail for delete", func() {
				err = client.StudentService.DeleteStudent(testCtx, missingId)
			})
		})
	})

	// FIXME: Break up Fun App Test logically and move to Handler Masterspec (Journal)
	Context("Admin", func() {
		It("should serve metrics", func() {
			err = client.AdminService.HealthCheck(testCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should serve swagger", func() {
			resp, err := DefaultHttpClient.R().Get(serviceUrl + "/swagger/index.html")

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(200))
		})
	})

})
