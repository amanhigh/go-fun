package play_test

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

type School struct {
	gorm.Model
	Name string
	City string
}

type Student struct {
	gorm.Model
	Name  string
	Age int
	School
}

type Teacher struct {
	gorm.Model
	Name  string
	Subject string
	School
}

var _ = FDescribe("Data Generator", Label(models.GINKGO_SETUP), func() {
	var (
		db  *gorm.DB
		err error
	)

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
