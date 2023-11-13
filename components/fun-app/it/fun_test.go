package it_test

import (
	"context"
	"strconv"
	"strings"

	. "github.com/amanhigh/go-fun/common/clients"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/fun"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Entr: http://eradman.com/entrproject/entr.1.html
// -s (use Shell), -c (Clear), Space/Q to Run, Quit.
// Watch Mode: find . | entr -s "date +%M:%S;ginkgo $PWD/components/fun-app/it | grep Pending"
var _ = Describe("Person Integration Test", func() {
	var (
		// serviceUrl = "http://localhost:8091/api/v1/namespaces/fun-app/services/fun-app:9000/proxy" //K8 endpoint or do PF on 8080 using K9S
		serviceUrl = "http://localhost:8085"
		request    fun.PersonRequest

		name   = "Amanpreet Singh"
		age    = 31
		gender = "MALE"
		client = NewFunAppClient(serviceUrl, config.DefaultHttpConfig)
		err    common.HttpError
		ctx    = context.Background()
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
			id string
		)
		BeforeEach(func() {
			id, err = client.PersonService.CreatePerson(ctx, request)
			Expect(id).Should(Not(BeEmpty()))
			Expect(err).To(BeNil())
		})
		AfterEach(func() {
			//Delete Person
			err = client.PersonService.DeletePerson(ctx, id)
			Expect(err).To(BeNil())
		})

		It("should create & get person", func() {
			person, err := client.PersonService.GetPerson(ctx, id)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(person).Should(Not(BeNil()))

			//Match Person Fields
			Expect(person.Id).To(Equal(id))
			Expect(person.Name).To(Equal(name))
			Expect(person.Age).To(Equal(age))
			Expect(person.Gender).To(Equal(gender))
		})

		Context("Update", func() {
			var (
				updateRequest fun.PersonRequest
				updateId      string
			)
			BeforeEach(func() {
				updateRequest = fun.PersonRequest{
					Name:   "Jenny",
					Age:    25,
					Gender: "FEMALE",
				}
				updateId, err = client.PersonService.CreatePerson(ctx, request)
				Expect(err).ShouldNot(HaveOccurred())
			})

			AfterEach(func() {
				err = client.PersonService.DeletePerson(ctx, updateId)
				Expect(err).To(BeNil())
			})

			It("should update person", func() {
				err := client.PersonService.UpdatePerson(ctx, updateId, updateRequest)
				Expect(err).ShouldNot(HaveOccurred())

				//Fetch Update Person
				person, err := client.PersonService.GetPerson(ctx, updateId)
				Expect(err).ShouldNot(HaveOccurred())

				//MatchFields
				Expect(person.Id).To(Equal(updateId))
				Expect(person.Name).To(Equal(updateRequest.Name))
				Expect(person.Age).To(Equal(updateRequest.Age))
				Expect(person.Gender).To(Equal(updateRequest.Gender))
			})

			Context("Bad Requests", func() {
				AfterEach(func() {
					err = client.PersonService.UpdatePerson(ctx, updateId, updateRequest)
					Expect(err).Should(HaveOccurred())
					Expect(err).To(Equal(common.ErrBadRequest))
				})

				It("should fail for missing Name", func() {
					updateRequest.Name = ""
				})

				It("should fail for invalid Name", func() {
					updateRequest.Name = "A*B"
				})

				It("should fail for max Name", func() {
					updateRequest.Name = strings.Repeat("A", 30)
				})

				It("should fail for missing Age", func() {
					updateRequest.Age = 0
				})

				It("should fail for invalid Age", func() {
					updateRequest.Age = -1
				})

				It("should fail for max Age", func() {
					updateRequest.Age = 200
				})

				It("should fail for missing Gender", func() {
					updateRequest.Gender = ""
				})

				It("should fail for invalid Gender", func() {
					updateRequest.Gender = "GENDER"
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
				Expect(personList.Total).To(BeEquivalentTo(total + 1))
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
				Expect(personList.Total).To(BeEquivalentTo(5))
			})

			It("should search by Gender", func() {
				var personList fun.PersonList
				personQuery.Gender = genders[1]
				personList, err = client.PersonService.ListPerson(ctx, personQuery)

				Expect(err).To(BeNil())
				Expect(len(personList.Records)).To(Equal(limit))
				Expect(personList.Total).To(BeEquivalentTo(11))
			})

			It("should search by Name & Gender", func() {
				var personList fun.PersonList
				personQuery.Name = names[0]
				personQuery.Gender = genders[1]
				personList, err = client.PersonService.ListPerson(ctx, personQuery)

				Expect(err).To(BeNil())
				Expect(len(personList.Records)).To(Equal(0))
				Expect(personList.Total).To(BeEquivalentTo(0))
			})

			Context("Bad Requests", func() {
				AfterEach(func() {
					_, err = client.PersonService.ListPerson(ctx, personQuery)
					Expect(err).Should(HaveOccurred())
					Expect(err).To(Equal(common.ErrBadRequest))
				})

				It("should fail for invalid Offset", func() {
					personQuery.Offset = -1
				})

				It("should fail for Lower Limit", func() {
					personQuery.Limit = 0
				})

				It("should fail for Max Limit", func() {
					personQuery.Limit = 30
				})

				It("should fail for invalid Name", func() {
					personQuery.Name = "A*B"
				})

				It("should fail for max Name", func() {
					personQuery.Name = strings.Repeat("A", 30)
				})

				It("should fail for invalid Gender", func() {
					personQuery.Gender = "OTHER"
				})
			})
		})

		Context("Bad Requests", func() {
			AfterEach(func() {
				_, err = client.PersonService.CreatePerson(ctx, request)

				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(common.ErrBadRequest))
				Expect(err.Code()).To(Equal(400))
			})

			It("should fail for missing Name", func() {
				request.Name = ""
			})

			It("should fail for invalid Name", func() {
				request.Name = "A*B"
			})

			It("should fail for max Name", func() {
				request.Name = strings.Repeat("A", 30)
			})

			It("should fail for minimum Age", func() {
				request.Age = 0
			})

			It("should fail for max Age", func() {
				request.Age = 200
			})

			It("should fail for missing Gender", func() {
				request.Gender = ""
			})

			It("should fail for invalid Gender", func() {
				request.Gender = "OTHER"
			})
		})
	})

	Context("Bad Requests", func() {
		var (
			emptyId   = ""
			missingId = "aba313bf"
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
			resp, err := TestHttpClient.R().
				Get(serviceUrl + "/swagger/index.html")

			Expect(err).To(BeNil())
			Expect(resp.StatusCode()).To(Equal(200))
		})
	})

})
