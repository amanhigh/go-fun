package it_test

import (
	"strings"

	. "github.com/amanhigh/go-fun/common/clients"
	"github.com/amanhigh/go-fun/models/common"
	db2 "github.com/amanhigh/go-fun/models/fun-app/db"
	server2 "github.com/amanhigh/go-fun/models/fun-app/server"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Person Integration Test", func() {
	var (
		// serviceUrl = "http://localhost:8091/api/v1/namespaces/fun-app/services/fun-app:9000/proxy" //K8 endpoint or do PF on 8080 using K9S
		serviceUrl = "http://localhost:8085"
		request    server2.PersonRequest

		name   = "Amanpreet Singh"
		age    = 31
		gender = "MALE"
		client = NewFunAppClient(serviceUrl)
		err    error
	)

	BeforeEach(func() {
		request = server2.PersonRequest{
			Person: db2.Person{
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
			// err = client.DeletePerson(name)
			// Expect(err).To(BeNil())
		})

		It("should create & get person", func() {
			person, err := client.PersonService.GetPerson(name)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(person).Should(Not(BeNil()))

			//Match Person Fields
			Expect(person.Id).ShouldNot(BeZero())
			Expect(person.Name).To(Equal(name))
			Expect(person.Age).To(Equal(age))
			Expect(person.Gender).To(Equal(gender))
		})

		It("should get all persons", func() {
			var persons []db2.Person
			persons, err = client.PersonService.GetAllPersons()
			Expect(err).To(BeNil())
			Expect(len(persons)).To(BeNumerically(">=", 1))
		})

		Context("Bad Requests", func() {
			AfterEach(func() {
				_, err = client.PersonService.CreatePerson(request)

				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(common.BadRequestErr))
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

	Context("Admin", func() {
		It("should serve metrics", func() {
			resp, err := TestHttpClient.R().
				Get(serviceUrl + "/metrics")

			Expect(err).To(BeNil())
			Expect(resp.StatusCode()).To(Equal(200))
		})

		It("should serve swagger", func() {
			resp, err := TestHttpClient.R().
				Get(serviceUrl + "/swagger/index.html")

			Expect(err).To(BeNil())
			Expect(resp.StatusCode()).To(Equal(200))
		})
	})

})
