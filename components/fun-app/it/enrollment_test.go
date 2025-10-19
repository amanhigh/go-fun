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
		var (
			initialRequest fun.EnrollmentRequest
		)

		Context("when grade is within capacity", func() {
			BeforeEach(func() {
				initialRequest = fun.EnrollmentRequest{PersonID: createdPerson.Id, Grade: 4}
				enrollResp, err = client.EnrollmentService.CreateEnrollment(ctx, initialRequest)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should create enrollment with ACTIVE status", func() {
				Expect(enrollResp.EnrollmentID).ToNot(BeEmpty())
				Expect(enrollResp.PersonID).To(Equal(createdPerson.Id))
				Expect(enrollResp.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
				Expect(enrollResp.Grade).To(Equal(initialRequest.Grade))

				getResp, getErr := client.EnrollmentService.GetEnrollment(ctx, createdPerson.Id)
				Expect(getErr).ToNot(HaveOccurred())
				Expect(getResp.EnrollmentID).To(Equal(enrollResp.EnrollmentID))
				Expect(getResp.PersonID).To(Equal(createdPerson.Id))
				Expect(getResp.Grade).To(Equal(initialRequest.Grade))
				Expect(getResp.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
			})

			Context("and enrollment is created again", func() {
				var (
					secondRequest fun.EnrollmentRequest
					secondResp    fun.EnrollmentResponse
					secondErr     common.HttpError
				)

				BeforeEach(func() {
					secondRequest = fun.EnrollmentRequest{PersonID: createdPerson.Id, Grade: 2}
					secondResp, secondErr = client.EnrollmentService.CreateEnrollment(ctx, secondRequest)
					Expect(secondErr).ToNot(HaveOccurred())
				})

				It("should update the existing enrollment", func() {
					Expect(secondResp.EnrollmentID).To(Equal(enrollResp.EnrollmentID))
					Expect(secondResp.Grade).To(Equal(secondRequest.Grade))
					Expect(secondResp.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))

					getResp, getErr := client.EnrollmentService.GetEnrollment(ctx, createdPerson.Id)
					Expect(getErr).ToNot(HaveOccurred())
					Expect(getResp.EnrollmentID).To(Equal(enrollResp.EnrollmentID))
					Expect(getResp.Grade).To(Equal(secondRequest.Grade))
					Expect(getResp.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
				})
			})
		})

		It("should waitlist when grade exceeds capacity", func() {
			enrollRequest = fun.EnrollmentRequest{PersonID: createdPerson.Id, Grade: 6}
			enrollResp, err = client.EnrollmentService.CreateEnrollment(ctx, enrollRequest)
			Expect(err).ToNot(HaveOccurred())
			Expect(enrollResp.Status).To(Equal(fun.EnrollmentStatusWaitlisted))

			getResp, getErr := client.EnrollmentService.GetEnrollment(ctx, createdPerson.Id)
			Expect(getErr).ToNot(HaveOccurred())
			Expect(getResp.Status).To(Equal(fun.EnrollmentStatusWaitlisted))
		})

		It("should return not found for unknown enrollment", func() {
			_, getErr := client.EnrollmentService.GetEnrollment(ctx, "missing-id")
			Expect(getErr).To(HaveOccurred())
			Expect(getErr.Code()).To(Equal(http.StatusNotFound))
		})
	})
})
