package play_test

import (
	"context"
	"math/rand"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/bxcodec/faker/v3"
	"github.com/caarlos0/env/v6"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

type School struct {
	gorm.Model `faker:"-"`
	Name       string `faker:"name" gorm:"not null,unique"`
	City       string `faker:"oneof:London,Berlin" gorm:"not null"`
}

type SchoolBase struct {
	School   School
	SchoolID uint `gorm:"not null"`
}

type Student struct {
	gorm.Model `faker:"-"`
	SchoolBase `faker:"-"`
	Name       string `faker:"name" gorm:"not null"`
	Birthday   int64  `faker:"unix_time" gorm:"not null"`
	Age        uint   `faker:"boundary_start=10,boundary_end=20"`
	Rank       uint   `faker:"oneof:1,2,3,4,5" gorm:"not null"`
}

type Teacher struct {
	gorm.Model `faker:"-"`
	SchoolBase `faker:"-"`
	Name       string `faker:"name" gorm:"not null"`
	Subject    string `faker:"oneof:Maths,Physics,Chemistry,English,History"`
	Phone      string `faker:"phone_number"`
	Salary     int    `faker:"boundary_start=20000,boundary_end=50000" gorm:"not null"`
}

var _ = Describe("Data Generator", Label(models.GINKGO_SETUP), func() {
	var (
		db  *gorm.DB
		ctx = context.Background()
		err error

		//Mysql Container
		container testcontainers.Container

		//Counts
		multiplier = 1
		batchSize  = 1000

		schoolCount  = multiplier * 3
		teacherCount = multiplier * 30
		studentCount = multiplier * 1000

		//Lists
		schools  = make([]School, schoolCount)
		students = make([]Student, studentCount)
		teachers = make([]Teacher, teacherCount)
	)

	BeforeEach(OncePerOrdered, func() {

		// Create School
		for i := 0; i < schoolCount; i++ {
			schools[i] = School{}
			err = faker.FakeData(&schools[i])
			Expect(err).To(BeNil())
		}

		// Create Teachers
		for i := 0; i < teacherCount; i++ {
			teachers[i] = Teacher{}
			err = faker.FakeData(&teachers[i])
			Expect(err).To(BeNil())
		}

		// Create Students
		for i := 0; i < studentCount; i++ {
			students[i] = Student{}
			err = faker.FakeData(&students[i])
			Expect(err).To(BeNil())
		}
	})

	It("should create required fake data", func() {
		Expect(len(schools)).To(Equal(schoolCount))
		Expect(len(students)).To(Equal(studentCount))
		Expect(len(teachers)).To(Equal(teacherCount))
	})

	Context("with db", Ordered, func() {
		var (
			dbconfig = config.Db{}
		)

		BeforeAll(func() {
			//Create Mysql Container
			container, err = util.MysqlTestContainer(ctx)
			Expect(err).To(BeNil())

			//Fill Defaults
			err = env.Parse(&dbconfig)
			Expect(err).To(BeNil())

			// url, err := container.Endpoint(ctx, "")
			mysqlPort, err := container.MappedPort(ctx, "3306")
			Expect(err).To(BeNil())

			dbconfig.Url = util.BuildMysqlURL("aman", "aman", "localhost", "compute", mysqlPort.Port())
			color.Green("Mysql Url: %s", dbconfig.Url)
			// dbconfig.LogLevel = logger.Info

			db, err = util.ConnectDb(dbconfig)
			Expect(err).To(BeNil())
		})

		AfterAll(func() {
			color.Red("Mysql Shutting Down")
			err = container.Terminate(ctx)
			Expect(err).To(BeNil())
		})

		It("should connect", func() {
			Expect(db).To(Not(BeNil()))
			Expect(err).To(BeNil())
		})

		Context("Migrate", func() {
			BeforeEach(func() {
				//Drop Existing Tables
				err = db.Migrator().DropTable(&School{}, &Student{}, &Teacher{})
				Expect(err).To(BeNil())

				err = db.AutoMigrate(&School{}, &Student{}, &Teacher{})
				Expect(err).To(BeNil())
			})

			It("should migrate", func() {
				Expect(db.Migrator().HasTable(&School{})).To(BeTrue(), "School")
				Expect(db.Migrator().HasTable(&Student{})).To(BeTrue(), "Student")
				Expect(db.Migrator().HasTable(&Teacher{})).To(BeTrue(), "Teacher")
			})

			Context("Insert", func() {
				var (
					actualCount int64
				)
				BeforeEach(func() {
					err = db.CreateInBatches(&schools, batchSize).Error
					Expect(err).To(BeNil())

					//Assign School Ids
					for teacher := range teachers {
						teachers[teacher].SchoolID = schools[rand.Intn(schoolCount)].ID
					}
					for student := range students {
						students[student].SchoolID = schools[rand.Intn(schoolCount)].ID
					}

					err = db.CreateInBatches(&teachers, batchSize).Error
					Expect(err).To(BeNil())

					err = db.CreateInBatches(&students, batchSize).Error
					Expect(err).To(BeNil())
				})

				It("should insert", func() {
					// Verify Insertion Count in DB
					err = db.Model(&School{}).Count(&actualCount).Error
					Expect(err).To(BeNil())
					Expect(actualCount).Should(BeNumerically("==", schoolCount))

					err = db.Model(&Student{}).Count(&actualCount).Error
					Expect(err).To(BeNil())
					Expect(actualCount).Should(BeNumerically("==", studentCount))

					err = db.Model(&Teacher{}).Count(&actualCount).Error
					Expect(err).To(BeNil())
					Expect(actualCount).Should(BeNumerically("==", teacherCount))
				})
			})
		})
	})
})
