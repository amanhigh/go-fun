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
	var db *gorm.DB

	BeforeEach(func() {
		db, _ = util.CreateTestDb(logger.Info)
	})

	It("should connect", func() {
		// Verify Connection
		Expect(db).To(Not(BeNil()))
	})

	Context("Create Table", func() {
		BeforeEach(func() {
			db.AutoMigrate(&frameworks.Product{}, &frameworks.AuditLog{})
		})

		AfterEach(func() {
			db.Migrator().DropTable(&frameworks.Product{})
			db.Migrator().DropTable(&frameworks.AuditLog{})
		})

		It("should create Tables", func() {
			// Verify Tables Created
			Expect(db.Migrator().HasTable(&frameworks.Product{})).To(BeTrue())
			Expect(db.Migrator().HasTable(&frameworks.AuditLog{})).To(BeTrue())
		})
	})

})
