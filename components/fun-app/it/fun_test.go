package it_test

import (
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/clients"
	"github.com/amanhigh/go-fun/models/fun"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FunApp Integration Smoke", func() {
	var (
		createdStudent fun.Student
		err            error

		retrievedStudent fun.Student
		updatedStudent   fun.Student
	)

	Context("Student lifecycle", func() {
		AfterEach(func() {
			if createdStudent.Id != "" {
				err = client.StudentService.DeleteStudent(ctx, createdStudent.Id)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		BeforeEach(func() {
			// 1. Create
			createdStudent, err = client.StudentService.CreateStudent(ctx, fun.StudentRequest{
				Name:   "Smoke Lifecycle",
				Age:    12,
				Gender: "MALE",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdStudent.Id).ToNot(BeEmpty())

			// 2. Retrieve
			retrievedStudent, err = client.StudentService.GetStudent(ctx, createdStudent.Id)
			Expect(err).ToNot(HaveOccurred())

			// 3. Update
			err = client.StudentService.UpdateStudent(ctx, createdStudent.Id, fun.StudentRequest{
				Name:   "Updated Smoke",
				Age:    20,
				Gender: "FEMALE",
			})
			Expect(err).ToNot(HaveOccurred())

			// 4. Retrieve updated
			updatedStudent, err = client.StudentService.GetStudent(ctx, createdStudent.Id)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should complete the student lifecycle", func() {
			Expect(retrievedStudent.Id).To(Equal(createdStudent.Id))
			Expect(retrievedStudent.Name).To(Equal("Smoke Lifecycle"))
			Expect(retrievedStudent.Age).To(Equal(12))
			Expect(retrievedStudent.Gender).To(Equal("MALE"))

			Expect(updatedStudent.Id).To(Equal(createdStudent.Id))
			Expect(updatedStudent.Name).To(Equal("Updated Smoke"))
			Expect(updatedStudent.Age).To(Equal(20))
			Expect(updatedStudent.Gender).To(Equal("FEMALE"))
		})

		Context("Enrollment lifecycle", func() {
			var (
				grade4Initial   fun.Enrollment
				grade4Confirmed fun.Enrollment
				grade6Initial   fun.Enrollment
				grade6Final     fun.Enrollment
			)

			BeforeEach(func() {
				// 1. Grade-4 enrollment
				grade4Req := fun.EnrollmentRequest{StudentID: createdStudent.Id, Grade: 4}
				grade4Initial, err = client.EnrollmentService.CreateEnrollment(ctx, grade4Req)
				Expect(err).ToNot(HaveOccurred())

				// 2. Poll until grade-4 enrollment is CONFIRMED
				Eventually(func() string {
					resp, pollErr := client.EnrollmentService.GetEnrollment(ctx, createdStudent.Id)
					if pollErr != nil {
						return ""
					}
					return resp.Status
				}, 5*time.Second, 100*time.Millisecond).Should(Equal(fun.EnrollmentStatusConfirmed))

				grade4Confirmed, err = client.EnrollmentService.GetEnrollment(ctx, createdStudent.Id)
				Expect(err).ToNot(HaveOccurred())

				// 3. Grade-6 re-enrollment
				grade6Req := fun.EnrollmentRequest{StudentID: createdStudent.Id, Grade: 6}
				grade6Initial, err = client.EnrollmentService.CreateEnrollment(ctx, grade6Req)
				Expect(err).ToNot(HaveOccurred())

				// 4. Poll checking both grade and status to avoid stale first-enrollment state
				Eventually(func() string {
					resp, pollErr := client.EnrollmentService.GetEnrollment(ctx, createdStudent.Id)
					if pollErr != nil {
						return ""
					}
					if resp.Grade != 6 {
						return ""
					}
					return resp.Status
				}, 5*time.Second, 100*time.Millisecond).Should(Equal(fun.EnrollmentStatusWaitlisted))

				grade6Final, err = client.EnrollmentService.GetEnrollment(ctx, createdStudent.Id)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should complete the enrollment lifecycle using the existing student", func() {
				// Grade-4 enrollment assertions
				Expect(grade4Initial.ID).ToNot(BeEmpty())
				Expect(grade4Initial.StudentID).To(Equal(createdStudent.Id))
				Expect(grade4Initial.Grade).To(Equal(4))
				Expect(grade4Initial.Status).To(Equal(fun.EnrollmentStatusSeatAllocationInitiated))

				Expect(grade4Confirmed.Grade).To(Equal(4))
				Expect(grade4Confirmed.Status).To(Equal(fun.EnrollmentStatusConfirmed))

				// Grade-6 re-enrollment assertions (same enrollment ID reused)
				Expect(grade6Initial.ID).To(Equal(grade4Initial.ID))
				Expect(grade6Initial.Grade).To(Equal(6))

				Expect(grade6Final.Grade).To(Equal(6))
				Expect(grade6Final.Status).To(Equal(fun.EnrollmentStatusWaitlisted))
			})
		})
	})

	Context("Operational endpoints", func() {
		var (
			metricsErr    error
			swaggerStatus int
		)

		BeforeEach(func() {
			metricsErr = client.AdminService.HealthCheck(ctx)

			resp, err := clients.DefaultHttpClient.R().
				SetContext(ctx).
				Get(serviceUrl + "/swagger/index.html")
			Expect(err).ToNot(HaveOccurred())
			swaggerStatus = resp.StatusCode()
		})

		It("should expose metrics and Swagger UI", func() {
			Expect(metricsErr).ToNot(HaveOccurred())
			Expect(swaggerStatus).To(Equal(http.StatusOK))
		})
	})
})
