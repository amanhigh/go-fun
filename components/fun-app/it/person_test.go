package it_test

import (
	. "github.com/amanhigh/go-fun/common/clients"
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
		BeforeEach(func() {
			err = client.CreatePerson(request)
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			//Delete Person
			// err = client.DeletePerson(name)
			// Expect(err).To(BeNil())
		})

		It("should create & get person", func() {
			person, err := client.GetPerson(name)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(person).Should(Not(BeNil()))

			//Match Person Fields
			Expect(person.Name).To(Equal(name))
			Expect(person.Age).To(Equal(age))
			Expect(person.Gender).To(Equal(gender))
		})

		It("should get all persons", func() {
			var persons []db2.Person
			persons, err = client.GetAllPersons()
			Expect(err).To(BeNil())
			Expect(len(persons)).To(BeNumerically(">=", 1))
		})

		PContext("Bad Requests", func() {
			It("should fail for bad request", func() {
				request.Name = ""
				err = client.CreatePerson(request)

				Expect(err).Should(HaveOccurred())
				//Error Should Contain Bad Request
				// Expect(err.Error()).To(ContainSubstring("Bad Request"))
			})
		})
	})

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
