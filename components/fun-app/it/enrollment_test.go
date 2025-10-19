package it_test

import (
	"net/http"
	"time"

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
		enrollResp    fun.Enrollment
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

			It("should create enrollment and confirm asynchronously", func() {
				Expect(enrollResp.ID).ToNot(BeEmpty())
				Expect(enrollResp.PersonID).To(Equal(createdPerson.Id))
				Expect(enrollResp.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
				Expect(enrollResp.Grade).To(Equal(initialRequest.Grade))

				Eventually(func() string {
					resp, pollErr := client.EnrollmentService.GetEnrollment(ctx, createdPerson.Id)
					if pollErr != nil {
						return ""
					}
					return resp.Status
				}, time.Second, 50*time.Millisecond).Should(Equal(fun.EnrollmentStatusConfirmed))

				getResp, getErr := client.EnrollmentService.GetEnrollment(ctx, createdPerson.Id)
				Expect(getErr).ToNot(HaveOccurred())
				Expect(getResp.ID).To(Equal(enrollResp.ID))
				Expect(getResp.PersonID).To(Equal(createdPerson.Id))
				Expect(getResp.Grade).To(Equal(initialRequest.Grade))
				Expect(getResp.Status).To(Equal(fun.EnrollmentStatusConfirmed))
			})

			Context("and enrollment is created again", func() {
				var (
					secondRequest fun.EnrollmentRequest
					secondResp    fun.Enrollment
					secondErr     common.HttpError
				)

				BeforeEach(func() {
					secondRequest = fun.EnrollmentRequest{PersonID: createdPerson.Id, Grade: 2}
					secondResp, secondErr = client.EnrollmentService.CreateEnrollment(ctx, secondRequest)
					Expect(secondErr).ToNot(HaveOccurred())
				})

				It("should update the existing enrollment", func() {
					Expect(secondResp.ID).To(Equal(enrollResp.ID))
					Expect(secondResp.Grade).To(Equal(secondRequest.Grade))
					Expect(secondResp.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))

					Eventually(func() string {
						resp, pollErr := client.EnrollmentService.GetEnrollment(ctx, createdPerson.Id)
						if pollErr != nil {
							return ""
						}
						if resp.Grade != secondRequest.Grade {
							return ""
						}
						return resp.Status
					}, time.Second, 50*time.Millisecond).Should(Equal(fun.EnrollmentStatusConfirmed))

					getResp, getErr := client.EnrollmentService.GetEnrollment(ctx, createdPerson.Id)
					Expect(getErr).ToNot(HaveOccurred())
					Expect(getResp.ID).To(Equal(enrollResp.ID))
					Expect(getResp.Grade).To(Equal(secondRequest.Grade))
					Expect(getResp.Status).To(Equal(fun.EnrollmentStatusConfirmed))
				})
			})
		})

		It("should waitlist when grade exceeds capacity", func() {
			enrollRequest = fun.EnrollmentRequest{PersonID: createdPerson.Id, Grade: 6}
			enrollResp, err = client.EnrollmentService.CreateEnrollment(ctx, enrollRequest)
			Expect(err).ToNot(HaveOccurred())
			Expect(enrollResp.Status).To(Equal(fun.EnrollmentStatusWaitlisted))

			Eventually(func() string {
				resp, pollErr := client.EnrollmentService.GetEnrollment(ctx, createdPerson.Id)
				if pollErr != nil {
					return ""
				}
				return resp.Status
			}, time.Second, 50*time.Millisecond).Should(Equal(fun.EnrollmentStatusWaitlisted))
		})

		It("should return not found for unknown enrollment", func() {
			_, getErr := client.EnrollmentService.GetEnrollment(ctx, "missing-id")
			Expect(getErr).To(HaveOccurred())
			Expect(getErr.Code()).To(Equal(http.StatusNotFound))
		})
	})
})
