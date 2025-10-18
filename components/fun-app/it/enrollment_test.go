package it_test

import (
	"net/http"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Enrollment API", func() {
	var (
		personRequest fun.PersonRequest
		createdPerson fun.Person
		enrollRequest fun.EnrollmentRequest
		enrollResp    fun.EnrollmentResponse
		err           common.HttpError
	)

	BeforeEach(func() {
		personRequest = fun.PersonRequest{
			Name:   "Saga Tester",
			Age:    10,
			Gender: "MALE",
		}

		createdPerson, err = client.PersonService.CreatePerson(ctx, personRequest)
		Expect(err).ToNot(HaveOccurred())
		Expect(createdPerson.Id).ToNot(BeEmpty())
	})

	AfterEach(func() {
		if createdPerson.Id != "" {
			err = client.PersonService.DeletePerson(ctx, createdPerson.Id)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Context("CreateEnrollment", func() {
		It("should enroll person when grade is within capacity", func() {
			enrollRequest = fun.EnrollmentRequest{PersonID: createdPerson.Id, Grade: 4}
			enrollResp, err = client.EnrollmentService.CreateEnrollment(ctx, enrollRequest)
			Expect(err).ToNot(HaveOccurred())
			Expect(enrollResp.PersonID).To(Equal(createdPerson.Id))
			Expect(enrollResp.Status).To(Equal("ACTIVE"))
			Expect(enrollResp.Grade).To(Equal(4))
		})

		It("should fail when grade exceeds capacity", func() {
			enrollRequest = fun.EnrollmentRequest{PersonID: createdPerson.Id, Grade: 6}
			_, err = client.EnrollmentService.CreateEnrollment(ctx, enrollRequest)
			Expect(err).To(HaveOccurred())
			Expect(err.Code()).To(Equal(http.StatusConflict))
		})
	})
})
