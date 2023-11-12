package it_test

import (
	"strconv"
	"strings"

	. "github.com/amanhigh/go-fun/common/clients"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun-app/db"
	"github.com/amanhigh/go-fun/models/fun-app/server"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Person Integration Test", func() {
	var (
		// serviceUrl = "http://localhost:8091/api/v1/namespaces/fun-app/services/fun-app:9000/proxy" //K8 endpoint or do PF on 8080 using K9S
		serviceUrl = "http://localhost:8085"
		request    server.PersonRequest

		name   = "Amanpreet Singh"
		age    = 31
		gender = "MALE"
		client = NewFunAppClient(serviceUrl)
		err    common.HttpError
	)

	BeforeEach(func() {
		request = server.PersonRequest{
			Person: db.Person{
				Name:   name,
				Age:    age,
				Gender: gender,
			},
		}
	})

	Context("Create", func() {
		var (
			id string
		)
		BeforeEach(func() {
			id, err = client.PersonService.CreatePerson(request)
			Expect(id).Should(Not(BeEmpty()))
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			//Delete Person
			err = client.PersonService.DeletePerson(id)
			Expect(err).To(BeNil())
		})

		It("should create & get person", func() {
			person, err := client.PersonService.GetPerson(id)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(person).Should(Not(BeNil()))

			//Match Person Fields
			Expect(person.Id).To(Equal(id))
			Expect(person.Name).To(Equal(name))
			Expect(person.Age).To(Equal(age))
			Expect(person.Gender).To(Equal(gender))
		})

		Context("Search", func() {
			var (
				offset      = 0
				limit       = 5
				total       = 15
				personQuery server.PersonQuery
				names       = []string{"Jane", "Sardar", "Rahul"}
				genders     = []string{"FEMALE", "MALE", "MALE"}
			)

			BeforeEach(func() {
				//Create 15 Persons
				for i := 0; i < total; i++ {
					request.Name = names[i%3] + strconv.Itoa(i)
					request.Gender = genders[i%3]
					_, err = client.PersonService.CreatePerson(request)
					Expect(err).To(BeNil())
				}

				//Init Person Query
				personQuery = server.PersonQuery{
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
					personList, err := client.PersonService.ListPerson(personQuery)
					Expect(err).To(BeNil())

					//Delete all Records of Name
					for _, person := range personList.Records {
						err = client.PersonService.DeletePerson(person.Id)
						Expect(err).To(BeNil())
					}
				}
			})

			It("should get all persons upto page Limit", func() {
				var personList server.PersonList
				personList, err = client.PersonService.ListPerson(personQuery)
				Expect(err).To(BeNil())

				//Person Count should be same as Page Limit
				Expect(len(personList.Records)).To(Equal(limit))
				Expect(personList.Total).To(BeEquivalentTo(total + 1))
			})

			It("should fetch second Page", func() {
				var personList server.PersonList
				personQuery.Offset = limit
				personList, err = client.PersonService.ListPerson(personQuery)

				Expect(err).To(BeNil())
				Expect(len(personList.Records)).To(Equal(limit))
			})

			It("should search by Name", func() {
				var personList server.PersonList
				personQuery.Name = names[0]
				personList, err = client.PersonService.ListPerson(personQuery)

				Expect(err).To(BeNil())
				Expect(len(personList.Records)).To(Equal(limit))
				Expect(personList.Total).To(BeEquivalentTo(5))
			})

			It("should search by Gender", func() {
				var personList server.PersonList
				personQuery.Gender = genders[1]
				personList, err = client.PersonService.ListPerson(personQuery)

				Expect(err).To(BeNil())
				Expect(len(personList.Records)).To(Equal(limit))
				Expect(personList.Total).To(BeEquivalentTo(11))
			})

			It("should search by Name & Gender", func() {
				var personList server.PersonList
				personQuery.Name = names[0]
				personQuery.Gender = genders[1]
				personList, err = client.PersonService.ListPerson(personQuery)

				Expect(err).To(BeNil())
				Expect(len(personList.Records)).To(Equal(0))
				Expect(personList.Total).To(BeEquivalentTo(0))
			})

			Context("Bad Requests", func() {
				AfterEach(func() {
					_, err = client.PersonService.ListPerson(personQuery)
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
				_, err = client.PersonService.CreatePerson(request)

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
				err = client.PersonService.DeletePerson(emptyId)
			})
		})

		Context("Missing Id", func() {
			AfterEach(func() {
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(common.ErrNotFound))
			})

			It("should fail for get", func() {
				_, err = client.PersonService.GetPerson(missingId)
			})

			It("should fail for delete", func() {
				err = client.PersonService.DeletePerson(missingId)
			})
		})
	})

	Context("Admin", func() {
		It("should serve metrics", func() {
			err = client.AdminService.HealthCheck()
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
