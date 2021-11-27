package it_test

import (
	"fmt"
	clients2 "github.com/amanhigh/go-fun/common/clients"
	db2 "github.com/amanhigh/go-fun/models/fun-app/db"
	server2 "github.com/amanhigh/go-fun/models/fun-app/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Person Integration Test", func() {
	var (
		serviceUrl = "http://localhost:8080"
		request    server2.PersonRequest

		name   = "Amanpreet Singh"
		age    = 31
		gender = "MALE"
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

		Context("Success", func() {
			AfterEach(func() {
				_, err := clients2.TestHttpClient.DoPost(serviceUrl+"/person", request, nil)
				Expect(err).To(BeNil())
			})

			It("should register person", func() {
			})
		})

		It("should fail for bad request", func() {
			request.Name = ""
			_, err := clients2.TestHttpClient.DoPost(serviceUrl+"/person", request, nil)
			Expect(err).To(Not(BeNil()))
		})
	})

	Context("Get", func() {

		It("should get all people", func() {
			var persons []db2.Person
			_, err := clients2.TestHttpClient.DoGet(serviceUrl+fmt.Sprintf("/person/all"), &persons)
			Expect(err).To(BeNil())
			Expect(len(persons)).To(BeNumerically(">=", 1))
		})

	})

	It("should serve metrics", func() {
		_, err := clients2.TestHttpClient.DoGet(serviceUrl+"/metrics", nil)
		Expect(err).To(BeNil())
	})
})
