package main_test

import (
	demo "github.com/amanhigh/go-fun/components/learn/frameworks/frontend/demo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StudentService", func() {
	It("seeds 20 students in ascending ID order", func() {
		service := demo.NewInMemoryStudentService()

		students, total := service.GetAllStudents(0, 20, "", "")

		Expect(students).To(HaveLen(20))
		Expect(total).To(Equal(20))
		Expect(students[0].ID).To(Equal("1"))
		Expect(students[19].ID).To(Equal("20"))
	})

	It("assigns the next ID to newly created students", func() {
		service := demo.NewInMemoryStudentService()

		created := service.CreateStudent(demo.Student{
			FirstName: "Zara",
			LastName:  "Khan",
			Email:     "zara.khan@school.edu",
			Age:       24,
			Grade:     "Senior",
		})

		students, total := service.GetAllStudents(0, 25, "", "")

		Expect(created.ID).To(Equal("21"))
		Expect(students).To(HaveLen(21))
		Expect(total).To(Equal(21))
		Expect(students[20].ID).To(Equal("21"))
	})

	It("filters by student name and grade before paginating", func() {
		service := demo.NewInMemoryStudentService()

		students, total := service.GetAllStudents(0, 10, "john", "Freshman")

		Expect(students).To(HaveLen(1))
		Expect(total).To(Equal(1))
		Expect(students[0].FirstName).To(Equal("Mike"))
		Expect(students[0].LastName).To(Equal("Johnson"))
	})
})
