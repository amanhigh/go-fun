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
			Expect(db.Migrator().HasTable(&frameworks.Vertical{})).To(BeTrue())
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
				err := db.FirstOrCreate(&vertical).Error
				Expect(err).To(BeNil())
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
				var vertical frameworks.Vertical
				err := db.Where(&frameworks.Vertical{Name: "Invalid"}).First(&vertical).Error
				Expect(err).To(Equal(gorm.ErrRecordNotFound))
			})

			Context("Create Product", func() {
				var product frameworks.Product

				BeforeEach(func() {
					features := []frameworks.Feature{
						{Name: "Strong", Version: 1},
						{Name: "Light", Version: 1},
					}
					product = frameworks.Product{Code: "L1212", Price: 1000, VerticalID: vertical.ID, Features: features, Version: 1}
					err := db.Create(&product).Error
					Expect(err).To(BeNil())
				})

				AfterEach(func() {
					db.Delete(&product)
					db.Exec("TRUNCATE TABLE features")
					db.Exec("TRUNCATE TABLE product_features")
				})

				It("should create product with features", func() {
					var foundProduct = new(frameworks.Product)
					err := db.Preload("Features").First(&foundProduct, product.ID).Error
					Expect(err).To(BeNil())
					Expect(foundProduct.Code).To(Equal("L1212"))
					Expect(foundProduct.Price).To(Equal(uint(1000)))
					Expect(foundProduct.Features).To(HaveLen(2))
				})

				It("should query product with code", func() {
					var queriedProduct frameworks.Product
					err := db.Preload("Vertical").First(&queriedProduct, "code = ?", product.Code).Error
					Expect(err).To(BeNil())
					Expect(queriedProduct.Vertical.Name).To(Equal(vertical.Name))
				})

				It("should query all non deleted products", func() {
					//Query all Non Deleted Products
					var products []frameworks.Product
					err := db.Unscoped().Where("code = ?", product.Code).Find(&products).Error
					Expect(err).To(BeNil())
					Expect(products).To(HaveLen(1))

				})

				It("should query id range", func() {
					//Query id range
					var products []frameworks.Product
					err := db.Where([]int64{5, 6, 10}).Limit(3).Limit(-1).Find(&products).Error
					Expect(err).To(BeNil())
					Expect(products).To(HaveLen(3))
				})

				It("should query multiple columns", func() {
					var multiSelectProducts []frameworks.Product
					err := db.Select("code", "price").Find(&multiSelectProducts).Error
					Expect(err).To(BeNil())
					Expect(multiSelectProducts).To(HaveLen(1))
					Expect(multiSelectProducts[0].Code).To(Equal(product.Code))
					Expect(multiSelectProducts[0].Price).To(Equal(product.Price))
					Expect(multiSelectProducts[0].VerticalID).To(Equal(uint(0))) // Not queried
				})

			})

		})
	})
})
