package it_test

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Enrollment API", func() {
	var (
		studentRequest fun.StudentRequest
		createdStudent fun.Student
		enrollRequest fun.EnrollmentRequest
		enrollResp    fun.Enrollment
		err           common.HttpError
	)

	BeforeEach(func() {
		studentRequest = fun.StudentRequest{
			Name:   "Saga Tester",
			Age:    10,
			Gender: "MALE",
		}

		createdStudent, err = client.StudentService.CreateStudent(ctx, studentRequest)
		Expect(err).ToNot(HaveOccurred())
		Expect(createdStudent.Id).ToNot(BeEmpty())
	})

	AfterEach(func() {
		if createdStudent.Id != "" {
			err = client.StudentService.DeleteStudent(ctx, createdStudent.Id)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Context("CreateEnrollment", func() {
		var (
			initialRequest fun.EnrollmentRequest
		)

		Context("when grade is within capacity", func() {
			BeforeEach(func() {
				initialRequest = fun.EnrollmentRequest{StudentID: createdStudent.Id, Grade: 4}
				enrollResp, err = client.EnrollmentService.CreateEnrollment(ctx, initialRequest)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should create enrollment and confirm asynchronously", func() {
				Expect(enrollResp.ID).ToNot(BeEmpty())
				Expect(enrollResp.StudentID).To(Equal(createdStudent.Id))
				Expect(enrollResp.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))
				Expect(enrollResp.Grade).To(Equal(initialRequest.Grade))

				Eventually(func() string {
					resp, pollErr := client.EnrollmentService.GetEnrollment(ctx, createdStudent.Id)
					if pollErr != nil {
						return ""
					}
					return resp.Status
				}, 3*time.Second, 50*time.Millisecond).Should(Equal(fun.EnrollmentStatusConfirmed))

				getResp, getErr := client.EnrollmentService.GetEnrollment(ctx, createdStudent.Id)
				Expect(getErr).ToNot(HaveOccurred())
				Expect(getResp.ID).To(Equal(enrollResp.ID))
				Expect(getResp.StudentID).To(Equal(createdStudent.Id))
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
					secondRequest = fun.EnrollmentRequest{StudentID: createdStudent.Id, Grade: 2}
					secondResp, secondErr = client.EnrollmentService.CreateEnrollment(ctx, secondRequest)
					Expect(secondErr).ToNot(HaveOccurred())
				})

				It("should update the existing enrollment", func() {
					Expect(secondResp.ID).To(Equal(enrollResp.ID))
					Expect(secondResp.Grade).To(Equal(secondRequest.Grade))
					Expect(secondResp.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))

					Eventually(func() string {
						resp, pollErr := client.EnrollmentService.GetEnrollment(ctx, createdStudent.Id)
						if pollErr != nil {
							return ""
						}
						if resp.Grade != secondRequest.Grade {
							return ""
						}
						return resp.Status
					}, 3*time.Second, 50*time.Millisecond).Should(Equal(fun.EnrollmentStatusConfirmed))

					getResp, getErr := client.EnrollmentService.GetEnrollment(ctx, createdStudent.Id)
					Expect(getErr).ToNot(HaveOccurred())
					Expect(getResp.ID).To(Equal(enrollResp.ID))
					Expect(getResp.Grade).To(Equal(secondRequest.Grade))
					Expect(getResp.Status).To(Equal(fun.EnrollmentStatusConfirmed))
				})
			})

			Context("and is re-enrolled above capacity", func() {
				var reEnrollResp fun.Enrollment

				BeforeEach(func() {
					Eventually(func() string {
						resp, pollErr := client.EnrollmentService.GetEnrollment(ctx, createdStudent.Id)
						if pollErr != nil {
							return ""
						}
						return resp.Status
					}, 3*time.Second, 50*time.Millisecond).Should(Equal(fun.EnrollmentStatusConfirmed))

					reEnrollResp, err = client.EnrollmentService.CreateEnrollment(ctx, fun.EnrollmentRequest{
						StudentID: createdStudent.Id,
						Grade:     6,
					})
					Expect(err).ToNot(HaveOccurred())
				})

				It("reuses the enrollment and waitlists the student", func() {
					Expect(reEnrollResp.ID).To(Equal(enrollResp.ID))
					Expect(reEnrollResp.StudentID).To(Equal(createdStudent.Id))
					Expect(reEnrollResp.Grade).To(Equal(6))
					Expect(reEnrollResp.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))

					Eventually(func() string {
						resp, pollErr := client.EnrollmentService.GetEnrollment(ctx, createdStudent.Id)
						if pollErr != nil || resp.Grade != 6 {
							return ""
						}
						return resp.Status
					}, 3*time.Second, 50*time.Millisecond).Should(Equal(fun.EnrollmentStatusWaitlisted))
				})
			})
		})

		It("should waitlist when grade exceeds capacity", func() {
			enrollRequest = fun.EnrollmentRequest{StudentID: createdStudent.Id, Grade: 6}
			enrollResp, err = client.EnrollmentService.CreateEnrollment(ctx, enrollRequest)
			Expect(err).ToNot(HaveOccurred())
			Expect(enrollResp.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))

			Eventually(func() string {
				resp, pollErr := client.EnrollmentService.GetEnrollment(ctx, createdStudent.Id)
				if pollErr != nil {
					return ""
				}
				return resp.Status
			}, time.Second, 50*time.Millisecond).Should(Equal(fun.EnrollmentStatusWaitlisted))
		})

	})
})
