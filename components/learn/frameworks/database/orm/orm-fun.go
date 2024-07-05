package orm

import (
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/models/learn/frameworks"
	. "github.com/amanhigh/go-fun/models/learn/frameworks"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func OrmFun() {
	//Can be Run Standalone for testing switch.
	//switchProduct()

	//schemaAlterPlay(db)
	// dropTables(db)
	fmt.Println("******ORM Fun Finished*******")
}

func switchProduct() {
	sourceCode := "Source Product"
	fmt.Println("***** Setting Up DB Resolver *****")

	db, err := gorm.Open(mysql.Open("aman:aman@tcp(mysql:3306)/compute?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{Logger: logger.Default.LogMode(1)})
	fmt.Println("Master DB Connect", err)

	/* Setup Resolver to Docker Mysql */
	db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{
			//sqlite.Open("/tmp/gorm.db"), //All Replica Calls Fail pointed to empty gorm.db since no replication
			mysql.Open("aman:aman@tcp(mysql:3307)/compute?charset=utf8&parseTime=True&loc=Local"),
		},
		Policy: dbresolver.RandomPolicy{},
	}))

	/* Migrate */
	db.AutoMigrate(&Product{})

	/* Write Source Products */
	vertical := frameworks.Vertical{
		Name:     "Test",
		MyColumn: "Hello",
	}
	db.FirstOrCreate(&vertical)
	//Auto Switch Writes to Source DB
	dbRes := db.FirstOrCreate(&Product{
		Code:     sourceCode,
		Price:    100,
		Version:  1,
		Vertical: vertical,
	})
	fmt.Println("Auto Switch Write: Write Success (Source)", dbRes.Error)
	fmt.Println("[Manual Switch and Write to Replica DB not Possible. Writes are forced to Sources.]")

	//Wait for Replication to Happen
	time.Sleep(2 * time.Second)

	/* Manual Switching Read */
	product := Product{}
	dbRes = db.Clauses(dbresolver.Write).Where("code = ?", sourceCode).Find(&product)
	fmt.Println("Manual Switch Read: Found (Source)", product.Code, len(product.Features), dbRes.Error)

	//BUG: Force Switch Preload Doesn't Work.
	//dbRes = db.Clauses(dbresolver.Write).Preload(clause.Associations, func(db *gorm.DB) *gorm.DB {
	//	return db.Clauses(dbresolver.Write)
	//}).Where("code = ?", sourceCode).Find(&product)

	product = Product{}
	dbRes = db.Clauses(dbresolver.Read).Where("code = ?", sourceCode).Find(&product)
	fmt.Println("Manual Switch Read: Found (Replica)", product.Code, dbRes.Error)

	/* Auto Switch Read */
	product = Product{}
	dbRes = db.Where("code = ?", sourceCode).Find(&product)
	fmt.Println("Auto Switch Read: Found (Replica)", product.Code, dbRes.Error)

	/* Transaction Read */
	product = Product{}
	err = db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("code = ?", sourceCode).Preload(clause.Associations).Find(&product).Error
	})
	fmt.Println("Transaction Read: Found (Source)", product.Code, len(product.Features), err)

}

func schemaAlterPlay(db *gorm.DB) {
	db.Migrator().DropColumn(&Product{}, "code")
}

func playProduct(db *gorm.DB) {
	fmt.Println("***** Play Product ******")
	productUpdates(db, product)
}

func productUpdates(db *gorm.DB, product *Product) {
	// Update without Callbacks
	db.Model(&product).UpdateColumn("code", "No Callback")
	//Single Field Update
	db.Model(&product).Update("Price", 1500)
	//Struct Update
	product.Code = "MyCode"
	db.Model(&product).Updates(product)

	manyToManyUpdate(db, product)

	fmt.Println("Product Updated")
}

func manyToManyUpdate(db *gorm.DB, product *Product) {
	fmt.Println("Before M2M Update: ", len(product.Features), product.Features[0].Name)

	//Existing Association needs to be deleted manually
	db.Delete(&product.Features[0])

	//New Associations can be added and Saved. Will not touch existing associations
	product.Features = []Feature{
		{Name: "abc", Version: 1},
		{Name: "xyz", Version: 1},
	}

	//Perform Update
	db.Save(product)

	//Only Displays New Features not the ones saved in DB's
	fmt.Println("Post M2M Update: ", len(product.Features), product.Features[0].Name)

	//Reload from Db
	reloadedProduct := Product{}
	db.Preload(clause.Associations).First(&reloadedProduct, product.ID)

	//Reloaded Product displays saved and newly created Features
	fmt.Println("Reloaded M2M Update (should have 3 features: 1 old, 2 replaced)", len(reloadedProduct.Features), reloadedProduct.Features[0].Name)
}
