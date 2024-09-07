package play_fast_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/learn/frameworks"
)

var _ = Describe("Orm", func() {
	var (
		db *gorm.DB
	)

	BeforeEach(func() {
		db, _ = util.CreateTestDb(logger.Error)
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
			db.Migrator().DropTable(&frameworks.Vertical{})
			db.Migrator().DropTable(&frameworks.Feature{})
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

		It("should drop the 'code' column from the Product table", func() {
			// Ensure column exists initially
			Expect(db.Migrator().HasColumn(&frameworks.Product{}, "code")).To(BeTrue())

			// Perform the column drop
			err := db.Migrator().DropColumn(&frameworks.Product{}, "code")
			Expect(err).ToNot(HaveOccurred())

			// Verify the column no longer exists
			Expect(db.Migrator().HasColumn(&frameworks.Product{}, "code")).To(BeFalse())
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
					err := db.Delete(&product).Error
					Expect(err).To(BeNil())

					// Existing Associations need to be deleted manually
					err = db.Exec("DELETE FROM features").Error
					Expect(err).To(BeNil())

					// Debug Features
					var features []frameworks.Feature
					err = db.Find(&features).Error
					Expect(err).To(BeNil())
					Expect(features).To(HaveLen(0))

					// Verify Feature Count
					var featureCount int64
					err = db.Model(&frameworks.Feature{}).Count(&featureCount).Error
					Expect(err).To(BeNil())
					Expect(featureCount).To(Equal(int64(0)))
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

				It("should query multiple columns", func() {
					var multiSelectProducts []frameworks.Product
					err := db.Select("code", "price").Find(&multiSelectProducts).Error
					Expect(err).To(BeNil())
					Expect(multiSelectProducts).To(HaveLen(1))
					Expect(multiSelectProducts[0].Code).To(Equal(product.Code))
					Expect(multiSelectProducts[0].Price).To(Equal(product.Price))
					Expect(multiSelectProducts[0].VerticalID).To(Equal(uint(0))) // Not queried
				})

				It("should pluck product codes", func() {
					var codes []string
					err := db.Model(&frameworks.Product{}).Pluck("code", &codes).Error
					Expect(err).To(BeNil())
					Expect(codes).To(HaveLen(1))
					Expect(codes[0]).To(Equal(product.Code))
				})

				Context("Bulk Products", func() {
					var products = []frameworks.Product{
						{Code: "L1212", Price: 1000, VerticalID: vertical.ID, Version: 1},
						{Code: "L1213", Price: 2000, VerticalID: vertical.ID, Version: 1},
						{Code: "L1214", Price: 3000, VerticalID: vertical.ID, Version: 1},
					}
					BeforeEach(func() {
						// Create multiple products for testing
						for _, product := range products {
							err := db.Create(&product).Error
							Expect(err).To(BeNil())
						}
					})
					AfterEach(func() {
						// Clean up the created products
						for _, product := range products {
							db.Delete(&product)
						}
					})

					It("should order query results", func() {
						var queriedProduct frameworks.Product
						err := db.Order("code desc, price asc").Last(&queriedProduct).Error
						Expect(err).To(BeNil())
						Expect(queriedProduct.Code).To(Equal(products[2].Code))
					})

					It("should query id range", func() {
						var products []frameworks.Product
						err := db.Where([]int64{1, 6, 10}).Limit(3).Limit(-1).Find(&products).Error
						Expect(err).To(BeNil())
						Expect(products).To(HaveLen(1))
					})

					It("should query by struct with OR condition", func() {
						var queriedProduct frameworks.Product
						err := db.Where(&frameworks.Product{Price: 2000}).Or(&frameworks.Product{Code: "Invalid"}).Last(&queriedProduct).Error
						Expect(err).To(BeNil())
						Expect(queriedProduct.Code).To(Equal(products[1].Code))
					})
				})

				Context("Update Product", func() {
					It("should update without callbacks", func() {
						err := db.Model(&product).UpdateColumn("code", "No Callback").Error
						Expect(err).To(BeNil())

						var updatedProduct frameworks.Product
						err = db.First(&updatedProduct, product.ID).Error
						Expect(err).To(BeNil())
						Expect(updatedProduct.Code).To(Equal("No Callback"))
						Expect(updatedProduct.Version).To(Equal(2)) // Version should not change
					})

					It("should update single field", func() {
						err := db.Model(&product).Update("Price", 1500).Error
						Expect(err).To(BeNil())

						var updatedProduct frameworks.Product
						err = db.First(&updatedProduct, product.ID).Error
						Expect(err).To(BeNil())
						Expect(updatedProduct.Price).To(Equal(uint(1500)))
						Expect(updatedProduct.Version).To(Equal(2)) // Version should increment
					})

					It("should update struct", func() {
						product.Code = "MyCode"
						err := db.Model(&product).Updates(product).Error
						Expect(err).To(BeNil())

						var updatedProduct frameworks.Product
						err = db.First(&updatedProduct, product.ID).Error
						Expect(err).To(BeNil())
						Expect(updatedProduct.Code).To(Equal("MyCode"))
						Expect(updatedProduct.Version).To(Equal(2)) // Version should increment
					})
				})

				Context("Many to Many Update", func() {
					var newFeatures = []frameworks.Feature{
						{Name: "abc", Version: 1},
						{Name: "xyz", Version: 1},
					}

					BeforeEach(func() {
						// Verify Current Features
						var initialFeatureCount int64
						err := db.Model(&frameworks.Feature{}).Count(&initialFeatureCount).Error
						Expect(err).To(BeNil())
						Expect(product.Features).To(HaveLen(int(initialFeatureCount)))
					})

					AfterEach(func() {
						// Delete new features
						for _, feature := range newFeatures {
							db.Where(frameworks.Feature{Name: feature.Name}).Delete(&feature)
						}
					})

					It("should delete 1 old, add 2 new", func() {
						db.Delete(&product.Features[0])

						// Add new associations
						product.Features = newFeatures

						// Perform Update
						err := db.Save(&product).Error
						Expect(err).To(BeNil())

						// Reload from DB
						var reloadedProduct frameworks.Product
						err = db.Preload(clause.Associations).First(&reloadedProduct, product.ID).Error
						Expect(err).To(BeNil())

						// // Check the updated features
						Expect(reloadedProduct.Features).To(HaveLen(3)) // 1 old + 2 new
						featureNames := []string{reloadedProduct.Features[0].Name, reloadedProduct.Features[1].Name, reloadedProduct.Features[2].Name}
						Expect(featureNames).To(ContainElements("Light", "abc", "xyz"))
					})

					It("should do association replacement", func() {
						err := db.Model(&product).Association("Features").Replace(newFeatures)
						Expect(err).To(BeNil())

						// Reload from DB
						var reloadedProduct frameworks.Product
						err = db.Preload("Features").First(&reloadedProduct, product.ID).Error
						Expect(err).To(BeNil())

						// Check that there are no features
						Expect(reloadedProduct.Features).To(HaveLen(2))
						Expect(reloadedProduct.Features[0].Name).To(Equal("abc"))
						Expect(reloadedProduct.Features[1].Name).To(Equal("xyz"))
					})
				})

			})

		})
	})
})
