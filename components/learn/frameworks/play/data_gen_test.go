package play_test

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/bxcodec/faker/v3"
	"github.com/caarlos0/env/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
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

var _ = Describe("Data Generator", Label(models.GINKGO_SLOW), func() {
	var (
		db  *gorm.DB
		ctx = context.Background()
		err error

		// Mysql Container
		container testcontainers.Container

		// Counts
		multiplier = 1
		batchSize  = 1000

		schoolCount  = multiplier * 3
		teacherCount = multiplier * 30
		studentCount = multiplier * 1000

		// Lists
		schools  = make([]School, schoolCount)
		students = make([]Student, studentCount)
		teachers = make([]Teacher, teacherCount)
	)

	BeforeEach(OncePerOrdered, func() {

		// Create School
		for i := 0; i < schoolCount; i++ {
			schools[i] = School{}
			err = faker.FakeData(&schools[i])
			Expect(err).ToNot(HaveOccurred())
		}

		// Create Teachers
		for i := 0; i < teacherCount; i++ {
			teachers[i] = Teacher{}
			err = faker.FakeData(&teachers[i])
			Expect(err).ToNot(HaveOccurred())
		}

		// Create Students
		for i := 0; i < studentCount; i++ {
			students[i] = Student{}
			err = faker.FakeData(&students[i])
			Expect(err).ToNot(HaveOccurred())
		}
	})

	It("should create required fake data", func() {

		Expect(schools).To(HaveLen(schoolCount))
		Expect(students).To(HaveLen(studentCount))
		Expect(teachers).To(HaveLen(teacherCount))
	})
	Context("with db", Ordered, func() {
		var (
			dbconfig = config.Db{}
		)

		BeforeAll(func() {
			// Create Mysql Container
			container, err = util.MysqlTestContainer(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Fill Defaults
			err = env.Parse(&dbconfig)
			Expect(err).ToNot(HaveOccurred())

			// url, err := container.Endpoint(ctx, "")
			mysqlPort, err := container.MappedPort(ctx, "3306")
			Expect(err).ToNot(HaveOccurred())

			dbconfig.DbType = models.MYSQL
			dbconfig.Url = util.BuildMysqlURL("aman", "aman", "localhost", "compute", mysqlPort.Port())
			log.Info().Str("URL", dbconfig.Url).Msg("Mysql Endpoint")
			// dbconfig.LogLevel = logger.Info

			db, err = util.ConnectDb(dbconfig)
			Expect(err).ToNot(HaveOccurred())
		})
		AfterAll(func() {
			log.Warn().Msg("Mysql Shutting Down")
			err = container.Terminate(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should connect", func() {
			Expect(db).To(Not(BeNil()))
			Expect(err).ToNot(HaveOccurred())
		})
		Context("Migrate", func() {
			BeforeEach(func() {
				// Drop Existing Tables
				err = db.Migrator().DropTable(&School{}, &Student{}, &Teacher{})
				Expect(err).ToNot(HaveOccurred())

				err = db.AutoMigrate(&School{}, &Student{}, &Teacher{})
				Expect(err).ToNot(HaveOccurred())
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
					Expect(err).ToNot(HaveOccurred())

					// Assign School Ids
					for teacher := range teachers {
						teachers[teacher].SchoolID = schools[util.RandomInt(0, schoolCount-1)].ID
					}
					for student := range students {
						students[student].SchoolID = schools[util.RandomInt(0, schoolCount-1)].ID
					}

					err = db.CreateInBatches(&teachers, batchSize).Error
					Expect(err).ToNot(HaveOccurred())

					err = db.CreateInBatches(&students, batchSize).Error
					Expect(err).ToNot(HaveOccurred())
				})
				It("should insert", func() {
					// Verify Insertion Count in DB
					err = db.Model(&School{}).Count(&actualCount).Error
					Expect(err).ToNot(HaveOccurred())
					Expect(actualCount).Should(BeNumerically("==", schoolCount))

					err = db.Model(&Student{}).Count(&actualCount).Error
					Expect(err).ToNot(HaveOccurred())
					Expect(actualCount).Should(BeNumerically("==", studentCount))
					err = db.Model(&Teacher{}).Count(&actualCount).Error
					Expect(err).ToNot(HaveOccurred())
					Expect(actualCount).Should(BeNumerically("==", teacherCount))
				})
			})
		})
	})
})
