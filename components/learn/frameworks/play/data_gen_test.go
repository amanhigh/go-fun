package play_test

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

type School struct {
	gorm.Model
	Name string `faker:"name" gorm:"unique,not null"`
	City string `faker:"oneof:London,Berlin" gorm:"not null"`
}

type Student struct {
	gorm.Model
	Name     string `faker:"name" gorm:"not null"`
	Birthday int64  `faker:"unix_time" gorm:"not null"`
	Age      uint   `faker:"boundary_start=10,boundary_end=20"`
	Rank     uint   `one_of:"1,2,3,4,5" gorm:"not null"`
	School
}

type Teacher struct {
	gorm.Model
	Name    string  `faker:"name" gorm:"not null"`
	Subject string  `one_of:"Maths,Physics,Chemistry,English,History"`
	Phone   string  `faker:"phone_number"`
	Salary  float64 `faker:"amount" gorm:"not null"`
	School
}

var _ = FDescribe("Data Generator", Label(models.GINKGO_SETUP), func() {
	var (
		db  *gorm.DB
		err error

		//Counts
		schoolCount  = 1
		studentCount = 10
		teacherCount = 5

		//Lists
		schools  = make([]School, schoolCount)
		students = make([]Student, studentCount)
		teachers = make([]Teacher, teacherCount)
	)

	BeforeEach(func() {
		// Create School
		for i := 0; i < schoolCount; i++ {
			schools[i] = School{}
			err = faker.FakeData(&schools[i])
			Expect(err).To(BeNil())
		}

		// Create Students
		for i := 0; i < studentCount; i++ {
			students[i] = Student{}
			err = faker.FakeData(&students[i])
			Expect(err).To(BeNil())
		}

		// Create Teachers
		for i := 0; i < teacherCount; i++ {
			teachers[i] = Teacher{}
			err = faker.FakeData(&teachers[i])
			Expect(err).To(BeNil())
		}
	})

	It("should create required fake data", func() {
		Expect(len(schools)).To(Equal(schoolCount))
		Expect(len(students)).To(Equal(studentCount))
		Expect(len(teachers)).To(Equal(teacherCount))
	})

	Context("with db", func() {
		BeforeEach(func() {
			db, err = util.CreateTestDb()
			Expect(err).To(BeNil())
		})

		It("should connect", func() {
			Expect(db).To(Not(BeNil()))
			Expect(err).To(BeNil())
		})

		Context("Migrate", func() {
			BeforeEach(func() {
				err = db.AutoMigrate(&School{}, &Student{}, &Teacher{})
				Expect(err).To(BeNil())
			})

			It("should migrate", func() {
				Expect(db.Migrator().HasTable(&School{})).To(BeTrue(), "School")
				Expect(db.Migrator().HasTable(&Student{})).To(BeTrue(), "Student")
				Expect(db.Migrator().HasTable(&Teacher{})).To(BeTrue(), "Teacher")
			})
		})
	})
})
