package it_test

import (
	"fmt"

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
			resp, err := TestHttpClient.R().
				SetBody(request).
				Post(serviceUrl + "/person")

			Expect(err).To(BeNil())
			Expect(resp.StatusCode()).To(Equal(200))
		})

		It("should create person", func() {
			person, err := client.GetPerson(name)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(person).Should(Not(BeNil()))
		})

		Context("Get", func() {
			It("should get all people", func() {
				var persons []db2.Person
				resp, err := TestHttpClient.R().
					SetResult(&persons).
					Get(serviceUrl + fmt.Sprintf("/person/all"))

				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(Equal(200))
				Expect(len(persons)).To(BeNumerically(">=", 1))
			})
		})

		It("should fail for bad request", func() {
			request.Name = ""
			resp, err := TestHttpClient.R().
				SetBody(request).
				Post(serviceUrl + "/person")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(400))
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
