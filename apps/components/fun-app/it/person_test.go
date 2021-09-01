package it_test

import (
	"fmt"

	"github.com/amanhigh/go-fun/apps/common/clients"
	"github.com/amanhigh/go-fun/apps/models/fun-app/db"
	"github.com/amanhigh/go-fun/apps/models/fun-app/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Person Integration Test", func() {
	var (
		serviceUrl = "http://localhost:8080"
		request    server.PersonRequest

		name   = "Amanpreet Singh"
		age    = 31
		gender = "MALE"
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

		Context("Success", func() {
			AfterEach(func() {
				_, err := clients.TestHttpClient.DoPost(serviceUrl+"/person", request, nil)
				Expect(err).To(BeNil())
			})

			It("should register person", func() {
			})
		})

		It("should fail for bad request", func() {
			request.Name = ""
			_, err := clients.TestHttpClient.DoPost(serviceUrl+"/person", request, nil)
			Expect(err).To(Not(BeNil()))
		})
	})

	Context("Get", func() {

		It("should get all people", func() {
			var persons []db.Person
			_, err := clients.TestHttpClient.DoGet(serviceUrl+fmt.Sprintf("/person/all"), &persons)
			Expect(err).To(BeNil())
			Expect(len(persons)).To(BeNumerically(">=", 1))
		})

	})

	It("should serve metrics", func() {
		_, err := clients.TestHttpClient.DoGet(serviceUrl+"/metrics", nil)
		Expect(err).To(BeNil())
	})
})
