package main_test

import (
	demo "github.com/amanhigh/go-fun/components/learn/frameworks/frontend/demo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StudentService", func() {
	It("seeds 20 students in ascending ID order", func() {
		service := demo.NewInMemoryStudentService()

		response := service.ListStudents(0, 20, "", "", "", "")

		Expect(response.Data).To(HaveLen(20))
		Expect(response.Count).To(Equal(20))
		Expect(response.Data[0].ID).To(Equal("1"))
		Expect(response.Data[19].ID).To(Equal("20"))
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

		response := service.ListStudents(0, 25, "", "", "", "")

		Expect(created.ID).To(Equal("21"))
		Expect(response.Data).To(HaveLen(21))
		Expect(response.Count).To(Equal(21))
		Expect(response.Data[20].ID).To(Equal("21"))
	})

	It("filters by student name and grade before paginating", func() {
		service := demo.NewInMemoryStudentService()

		response := service.ListStudents(0, 10, "john", "Freshman", "", "")

		Expect(response.Data).To(HaveLen(1))
		Expect(response.Count).To(Equal(1))
		Expect(response.Data[0].FirstName).To(Equal("Mike"))
		Expect(response.Data[0].LastName).To(Equal("Johnson"))
	})

	It("sorts by name when requested", func() {
		service := demo.NewInMemoryStudentService()

		response := service.ListStudents(0, 5, "", "", "name", "desc")

		Expect(response.Data).To(HaveLen(5))
		Expect(response.Data[0].FirstName).To(Equal("Sophia"))
		Expect(response.Data[0].LastName).To(Equal("Thomas"))
	})
})
