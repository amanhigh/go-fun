package play_fast_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/learn/frameworks"
)

var _ = FDescribe("Orm", func() {
	var (
		db *gorm.DB
	)

	BeforeEach(func() {
		db, _ = util.CreateTestDb(logger.Info)
	})

	It("should connect", func() {
		// Verify Connection
		Expect(db).To(Not(BeNil()))
	})

	Context("Create Table", func() {
		BeforeEach(func() {
			err := db.AutoMigrate(&frameworks.Product{}, &frameworks.AuditLog{})
			Expect(err).To(BeNil())
			db.Commit()
		})

		AfterEach(func() {
			db.Migrator().DropTable(&frameworks.Product{})
			err := db.Migrator().DropTable(&frameworks.AuditLog{})
			Expect(err).To(BeNil())
		})

		It("should create Tables", func() {
			// Verify Tables Created
			Expect(db.Migrator().HasTable(&frameworks.Product{})).To(BeTrue())
			Expect(db.Migrator().HasTable(&frameworks.AuditLog{})).To(BeTrue())
		})

		It("should have no records", func() {
			// Verify Tables Created
			var count int64
			err := db.Model(&frameworks.Product{}).Count(&count).Error
			Expect(err).To(BeNil())
			Expect(count).To(Equal(int64(0)))
		})

		Context("Create Vertical", func() {
			var vertical frameworks.Vertical

			BeforeEach(func() {
				vertical = frameworks.Vertical{
					Name:     "Test",
					MyColumn: "Hello",
				}
				err := db.FirstOrCreate(&vertical)
				Expect(err).To(BeNil())
				db.Commit()
			})

			AfterEach(func() {
				db.Exec("truncate table verticals")
			})

			It("should have create vertical", func() {
				vertical := frameworks.Vertical{}
				err := db.First(&vertical).Error
				Expect(err).To(BeNil())
				Expect(vertical.Name).To(Equal("Test"))
				Expect(vertical.MyColumn).To(Equal("Hello"))
			})

			It("should not find invalid vertical", func() {
				err := db.Model(&frameworks.Vertical{Name: "Invalid"}).Error
				Expect(err).To(BeNil())
				Expect(err).To(Equal(gorm.ErrRecordNotFound))
			})
		})
	})
})
