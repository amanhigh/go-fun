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
var _ = Describe("Person Integration Test", func() {
	var (
		request fun.PersonRequest

		name        = "Amanpreet Singh"
		maxName     = strings.Repeat("A", 31)
		age         = 31
		gender      = "MALE"
		err         common.HttpError
		ctx         = context.Background()
		expectedErr = "Bad Request"
	)

	BeforeEach(func() {
		request = fun.PersonRequest{
			Name:   name,
			Age:    age,
			Gender: gender,
		}
	})

	Context("Create", func() {
		var (
			createdPerson fun.Person
			auditUser     = "AMAN"
		)
		BeforeEach(func() {
			createdPerson, err = client.PersonService.CreatePerson(ctx, request)
			Expect(err).To(BeNil())
			Expect(createdPerson.Id).Should(Not(BeEmpty()))
		})
		AfterEach(func() {
			//Delete Person
			err = client.PersonService.DeletePerson(ctx, createdPerson.Id)
			Expect(err).To(BeNil())

			//Delete Audit
			auditList, err := client.PersonService.ListPersonAudit(ctx, createdPerson.Id)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(auditList)).To(Equal(2))
		})

		It("should create & get person", func() {
			person, err := client.PersonService.GetPerson(ctx, createdPerson.Id)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(person).Should(Not(BeNil()))

			//Match Person Fields
			Expect(person.Id).To(Equal(createdPerson.Id))
			Expect(person.Name).To(Equal(name))
			Expect(person.Age).To(Equal(age))
			Expect(person.Gender).To(Equal(gender))
		})

		It("should generate Audit", func() {
			//List Audit
			auditList, err := client.PersonService.ListPersonAudit(ctx, createdPerson.Id)
			Expect(err).ShouldNot(HaveOccurred())

			//Check Audit
			Expect(len(auditList)).To(Equal(1))
			audit := auditList[0]
			Expect(audit.Id).To(Equal(createdPerson.Id))
			Expect(audit.Name).To(Equal(name))
			Expect(audit.Age).To(Equal(age))
			Expect(audit.Gender).To(Equal(gender))

			Expect(audit.Operation).To(Equal("CREATE"))
			Expect(audit.CreatedBy).To(Equal(auditUser))
			Expect(audit.CreatedAt).Should(Not(BeNil()))
		})

		Context("Update", func() {
			var (
				updateRequest fun.PersonRequest
				updatedPerson fun.Person
			)
			BeforeEach(func() {
				updateRequest = fun.PersonRequest{
					Name:   "Jenny",
					Age:    25,
					Gender: "FEMALE",
				}
				updatedPerson, err = client.PersonService.CreatePerson(ctx, request)
				Expect(err).ShouldNot(HaveOccurred())
			})

			AfterEach(func() {
				err = client.PersonService.DeletePerson(ctx, updatedPerson.Id)
				Expect(err).To(BeNil())
			})

			Context("Success", func() {
				BeforeEach(func() {
					err := client.PersonService.UpdatePerson(ctx, updatedPerson.Id, updateRequest)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("should update person", func() {
					//Fetch Update Person
					person, err := client.PersonService.GetPerson(ctx, updatedPerson.Id)
					Expect(err).ShouldNot(HaveOccurred())

					//MatchFields
					Expect(person.Id).To(Equal(updatedPerson.Id))
					Expect(person.Name).To(Equal(updateRequest.Name))
					Expect(person.Age).To(Equal(updateRequest.Age))
					Expect(person.Gender).To(Equal(updateRequest.Gender))
				})

				It("should generate Audit", func() {
					//List Audit
					auditList, err := client.PersonService.ListPersonAudit(ctx, updatedPerson.Id)
					Expect(err).ShouldNot(HaveOccurred())

					//Check Audit
					Expect(len(auditList)).To(Equal(2))
					audit := auditList[1]
					Expect(audit.Id).To(Equal(updatedPerson.Id))
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
					err = client.PersonService.UpdatePerson(ctx, updatedPerson.Id, updateRequest)
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(expectedErr))
					Expect(err.Code()).To(Equal(http.StatusBadRequest))
				})

				It("should fail for missing Name", func() {
					updateRequest.Name = ""
					expectedErr = "required"
				})

				It("should fail for invalid Name", func() {
					updateRequest.Name = "A*B"
					expectedErr = "Name"
				})

				It("should fail for max Name", func() {
					updateRequest.Name = maxName
					expectedErr = "max"
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
					expectedErr = "max"
				})

				It("should fail for missing Gender", func() {
					updateRequest.Gender = ""
					expectedErr = "required"
				})

				It("should fail for invalid Gender", func() {
					updateRequest.Gender = "GENDER"
					expectedErr = "FEMALE"
				})

			})
		})

		Context("Search", func() {
			var (
				offset      = 0
				limit       = 5
				total       = 15
				personQuery fun.PersonQuery
				names       = []string{"Jane", "Sardar", "Rahul"}
				genders     = []string{"FEMALE", "MALE", "MALE"}
			)

			BeforeEach(func() {
				//Create 15 Persons
				for i := 0; i < total; i++ {
					request.Name = names[i%3] + strconv.Itoa(i)
					request.Gender = genders[i%3]
					_, err = client.PersonService.CreatePerson(ctx, request)
					Expect(err).To(BeNil())
				}

				//Init Person Query
				personQuery = fun.PersonQuery{
					Pagination: common.Pagination{
						Offset: offset,
						Limit:  limit,
					},
				}
			})

			AfterEach(func() {
				//Find Record By Names and Delete using UUID
				for i, name := range names {
					personQuery.Name = name
					personQuery.Gender = genders[i]
					personQuery.Limit = 10
					personQuery.Offset = 0
					personList, err := client.PersonService.ListPerson(ctx, personQuery)
					Expect(err).To(BeNil())

					//Delete all Records of Name
					for _, person := range personList.Records {
						err = client.PersonService.DeletePerson(ctx, person.Id)
						Expect(err).To(BeNil())
					}
				}
			})

			It("should get all persons upto page Limit", func() {
				var personList fun.PersonList
				personList, err = client.PersonService.ListPerson(ctx, personQuery)
				Expect(err).To(BeNil())

				//Person Count should be same as Page Limit
				Expect(len(personList.Records)).To(Equal(limit))
				Expect(personList.Metadata.Total).To(BeNumerically(">=", total))
			})

			It("should fetch second Page", func() {
				var personList fun.PersonList
				personQuery.Offset = limit
				personList, err = client.PersonService.ListPerson(ctx, personQuery)

				Expect(err).To(BeNil())
				Expect(len(personList.Records)).To(Equal(limit))
			})

			It("should search by Name", func() {
				var personList fun.PersonList
				personQuery.Name = names[0]
				personList, err = client.PersonService.ListPerson(ctx, personQuery)

				Expect(err).To(BeNil())
				Expect(len(personList.Records)).To(Equal(limit))
				Expect(personList.Metadata.Total).To(BeEquivalentTo(5))
			})

			It("should search by Gender", func() {
				var personList fun.PersonList
				personQuery.Gender = genders[1]
				personList, err = client.PersonService.ListPerson(ctx, personQuery)

				Expect(err).To(BeNil())
				Expect(len(personList.Records)).To(Equal(limit))
				Expect(personList.Metadata.Total).To(BeEquivalentTo(11))
			})

			It("should search by Name & Gender", func() {
				var personList fun.PersonList
				personQuery.Name = names[0]
				personQuery.Gender = genders[1]
				personList, err = client.PersonService.ListPerson(ctx, personQuery)

				Expect(err).To(BeNil())
				Expect(len(personList.Records)).To(Equal(0))
				Expect(personList.Metadata.Total).To(BeEquivalentTo(0))
			})

			Context("Sort", func() {

				It("should sort by Name in ascending order", func() {
					var personList fun.PersonList
					personQuery.SortBy = "name"
					personQuery.Order = "asc"
					personList, err = client.PersonService.ListPerson(ctx, personQuery)

					Expect(err).To(BeNil())
					Expect(len(personList.Records)).To(Equal(limit))
					// Check if the records are sorted in ascending order by name
					for i := 0; i < len(personList.Records)-1; i++ {
						cur := personList.Records[i].Name
						next := personList.Records[i+1].Name
						Expect(cur <= next).To(BeTrue())
					}
				})

				It("should sort by Name in descending order", func() {
					var personList fun.PersonList
					personQuery.SortBy = "name"
					personQuery.Order = "desc"
					personList, err = client.PersonService.ListPerson(ctx, personQuery)

					Expect(err).To(BeNil())
					Expect(len(personList.Records)).To(Equal(limit))

					// Check if the records are sorted in descending order by name
					for i := 0; i < len(personList.Records)-1; i++ {
						cur := personList.Records[i].Name
						next := personList.Records[i+1].Name
						Expect(cur >= next).To(BeTrue())
					}
				})

				It("should sort by Gender in ascending order", func() {
					var personList fun.PersonList
					personQuery.SortBy = "gender"
					personQuery.Order = "asc"
					personList, err = client.PersonService.ListPerson(ctx, personQuery)

					Expect(err).To(BeNil())
					Expect(len(personList.Records)).To(Equal(limit))

					// Check if the records are sorted in ascending order by gender
					for i := 0; i < len(personList.Records)-1; i++ {
						cur := personList.Records[i].Gender
						next := personList.Records[i+1].Gender
						Expect(cur <= next).To(BeTrue())
					}
				})

				It("should sort by Gender in descending order", func() {
					var personList fun.PersonList
					personQuery.SortBy = "gender"
					personQuery.Order = "desc"
					personList, err = client.PersonService.ListPerson(ctx, personQuery)

					Expect(err).To(BeNil())
					Expect(len(personList.Records)).To(Equal(limit))

					// Check if the records are sorted in descending order by gender
					for i := 0; i < len(personList.Records)-1; i++ {
						cur := personList.Records[i].Gender
						next := personList.Records[i+1].Gender
						Expect(cur >= next).To(BeTrue())
					}
				})
			})

			Context("Bad Requests", func() {
				AfterEach(func() {
					_, err = client.PersonService.ListPerson(ctx, personQuery)
					Expect(err).Should(HaveOccurred())
					Expect(err.Code()).To(Equal(http.StatusBadRequest))
					Expect(err.Error()).To(ContainSubstring(expectedErr))

					//Pollutes AfterEach Cleanup so Reset
					personQuery.Order = ""
					personQuery.SortBy = ""
				})

				It("should fail for invalid Offset", func() {
					personQuery.Offset = -1
					expectedErr = "Offset"
				})

				It("should fail for Lower Limit", func() {
					personQuery.Limit = 0
					expectedErr = "required"
				})

				It("should fail for Max Limit", func() {
					personQuery.Limit = 30
					expectedErr = "Limit"
				})

				It("should fail for invalid Name", func() {
					personQuery.Name = "A*B"
					expectedErr = "Name"
				})

				It("should fail for max Name", func() {
					personQuery.Name = maxName
					expectedErr = "max"
				})

				It("should fail for invalid Gender", func() {
					personQuery.Gender = "OTHER"
					expectedErr = "FEMALE"
				})

				It("should fail for invalid SortBy", func() {
					personQuery.SortBy = "invalid"
					expectedErr = "SortBy"
				})

				It("should fail for invalid Order", func() {
					personQuery.SortBy = "name"
					personQuery.Order = "invalid"
					expectedErr = "asc"
				})
			})
		})

		Context("Bad Requests", func() {
			AfterEach(func() {
				_, err = client.PersonService.CreatePerson(ctx, request)

				Expect(err).Should(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusBadRequest))
				Expect(err.Error()).To(ContainSubstring(expectedErr))
			})

			It("should fail for missing Name", func() {
				request.Name = ""
				expectedErr = "required"
			})

			It("should fail for invalid Name", func() {
				request.Name = "A*B"
				expectedErr = "Name"
			})

			It("should fail for max Name", func() {
				request.Name = maxName
				expectedErr = "max"
			})

			It("should fail for minimum Age", func() {
				request.Age = 0
				expectedErr = "Age"
			})

			It("should fail for max Age", func() {
				request.Age = 200
				expectedErr = "max"
			})

			It("should fail for missing Gender", func() {
				request.Gender = ""
				expectedErr = "Gender"
			})

			It("should fail for invalid Gender", func() {
				request.Gender = "OTHER"
				expectedErr = "FEMALE"
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
				err = client.PersonService.DeletePerson(ctx, emptyId)
			})
		})

		Context("Missing Id", func() {
			AfterEach(func() {
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(common.ErrNotFound))
			})

			It("should fail for get", func() {
				_, err = client.PersonService.GetPerson(ctx, missingId)
			})

			It("should fail for delete", func() {
				err = client.PersonService.DeletePerson(ctx, missingId)
			})
		})
	})

	Context("Admin", func() {
		It("should serve metrics", func() {
			err = client.AdminService.HealthCheck(ctx)
			Expect(err).To(BeNil())
		})

		It("should serve swagger", func() {
			resp, err := TestHttpClient.R().Get(serviceUrl + "/swagger/index.html")

			Expect(err).To(BeNil())
			Expect(resp.StatusCode()).To(Equal(200))
		})
	})

})
