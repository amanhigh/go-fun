package orm

import (
	"encoding/json"
	"errors"
	"fmt"
	util2 "github.com/amanhigh/go-fun/apps/common/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"

	"github.com/amanhigh/go-fun/learn/frameworks/orm/model"
	_ "github.com/amanhigh/go-fun/util"
	log "github.com/sirupsen/logrus"
)

type Product struct {
	gorm.Model
	Code       string `gorm:"size 5"`
	Price      uint
	Version    int
	IgnoreMe   string `gorm:"-"` // Ignore this field
	Vertical   model.Vertical
	VerticalID uint      //Must be vertical_id in DB or won't work automatically.
	Features   []Feature `gorm:"many2many:product_features;"`
}

type AuditLog struct {
	gorm.Model
	Operation string
	Log       string
}

type Feature struct {
	gorm.Model
	Name string
}

//Default Name would be products
func (p *Product) TableName() string {
	return "MeraProduct"
}

func (u *Product) BeforeCreate(tx *gorm.DB) (err error) {
	//Backup Product
	marshal, _ := json.Marshal(u)
	u.Version += 1
	tx.Create(&AuditLog{Operation: "Create", Log: string(marshal)})
	return
}

func (u *Product) BeforeUpdate(tx *gorm.DB) (err error) {
	//Backup Product
	marshal, _ := json.Marshal(u)
	u.Version += 1
	tx.Create(&AuditLog{Operation: "Update", Log: string(marshal)})
	return
}

// begin transaction
// -> BeforeSave
// -> BeforeCreate/Update
// save before associations
// update timestamp `CreatedAt`, `UpdatedAt`
// save self
// reload fields that have default value and its value is blank
// save after associations
// -> AfterCreate
// -> AfterSave/Update
// commit or rollback transaction

func (p *Product) AfterFind(_ *gorm.DB) (err error) {
	p.IgnoreMe = "Ignore" + p.Code
	return nil
}

/** Extra Json Logger */
var jsonLogger = &log.Logger{Out: os.Stdout, Formatter: new(log.JSONFormatter), Level: log.InfoLevel}

func OrmFun() {
	db, _ := util2.CreateTestDb()

	prepLogger()
	db.AutoMigrate(&Product{}, &AuditLog{}) // Vertical not required Foreign Keys Auto Created

	playProduct(db)

	//schemaAlterPlay(db)
	fmt.Println("ORM Fun Finished")
}

func prepLogger() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func schemaAlterPlay(db *gorm.DB) {
	db.Migrator().DropColumn(&Product{}, "code")
}

func TruncateTable(db *gorm.DB, tableName string) {
	db.Exec("truncate table " + tableName)
}

func playProduct(db *gorm.DB) {
	createVertical(db)

	// Create
	features := []Feature{
		{Name: "Strong"},
		{Name: "Light"},
	}
	product := &Product{Code: "L1212", Price: 1000, VerticalID: 1, Features: features, Version: 1}
	db.Create(product)

	queryProduct(db)

	productUpdates(db, product)

	productBackup(db)

	// Delete - delete product
	db.Delete(&product)
}

func productBackup(db *gorm.DB) {
	var count int64
	db.Model(&AuditLog{}).Count(&count)
	fmt.Println("Audit Logs: ", count)
}
func productUpdates(db *gorm.DB, product *Product) {
	// Update without Callbacks
	db.Model(&product).UpdateColumn("code", "No Callback")
	//Single Field Update
	db.Model(&product).Update("Price", 1500)
	//Struct Update
	db.Model(&product).Updates(&Product{Code: "MyCode"})

	manyToManyUpdate(db, product)

	fmt.Println("Product Updated")
}

func manyToManyUpdate(db *gorm.DB, product *Product) {
	fmt.Println("Before M2M Update: ", len(product.Features), product.Features[0].Name)

	//Existing Assocation needs to be deleted manually
	db.Delete(product.Features[0])

	//New Associations can be added and Saved. Will not touch existing associations
	product.Features = []Feature{
		{Name: "abc"},
		{Name: "xyz"},
	}

	//Perform Update
	db.Save(product)

	//Only Displays New Features not the ones saved in DB's
	fmt.Println("Post M2M Update: ", len(product.Features), product.Features[0].Name)

	//Reload from Db
	reloadedProduct := Product{}
	db.Preload(clause.Associations).First(&reloadedProduct, product.ID)

	//Reloaded Product displays saved and newly created Features
	//TODO: Deleted Items are coming why ?
	fmt.Println("Reloaded M2M Update", len(reloadedProduct.Features), reloadedProduct.Features[0].Name)
}

func queryProduct(db *gorm.DB) {
	// First Query
	product := new(Product)
	db.First(product, 1)

	// Preload with Where Clause
	db.Preload("Vertical").First(product, "code = ?", "L1212")

	fields := log.Fields{
		"Vertical ID:":  product.VerticalID,
		"Vertical Name": product.Vertical.Name,
		"Ignore Me":     product.IgnoreMe,
	}
	log.WithFields(fields).Info("Product Details")
	jsonLogger.WithFields(fields).Info("Product Details")

	//Query all Non Deleted Products
	products := new([]Product)
	db.Unscoped().Where("code = ?", "L1212").Find(products)
	for _, product := range *products {
		fmt.Println("Deleted/Undeleted Product Found: ", product.ID)
	}

	//Single Field Select
	//TODO:Fix
	//var codes []string
	//db.Not([]int64{5, 6, 10}).Find(products).Pluck("code", &codes)
	//fmt.Printf("CODES: %+v\n", codes)

	//Multi Field Select
	var multiSelectProducts []Product
	db.Select("code", "price").Find(&multiSelectProducts)
	fmt.Println("Multi Field Select")
	for _, p := range multiSelectProducts {
		//Vertical Id is not queried
		fmt.Println(p.Code, p.Price, p.VerticalID)
	}

	//Search Id Range
	db.Unscoped().Where([]int64{5, 6, 10}).Limit(3).Limit(-1).Find(products)
	fmt.Println("Id Range Search Count: ", len(*products))

	//Struct Query
	db.Order("code desc,price asc").Where(&Product{Price: 2000}).Where(&Product{Code: "L1212"}).Last(product) //And
	db.Where(&Product{Price: 2000}).Or(&Product{Code: "L1212"}).Last(product)                                 //Or
	fmt.Println("Query By Struct, ID:", product.ID)
}

func createVertical(db *gorm.DB) {
	vertical := &model.Vertical{}
	db.FirstOrCreate(&vertical)
	count := new(int64)
	db.Model(&model.Vertical{}).Count(count)
	fmt.Println("Vertical Count:", *count)

	fmt.Println("\n\nVertical Json WRITE")
	vertical.WriteTo(os.Stdout)
	fmt.Println("\nVertical Json WRITE\n\n")

	result := db.Model(&model.Vertical{Name: "Not Present"})
	fmt.Println("FOUND Value:", errors.Is(result.Error, gorm.ErrRecordNotFound))

	//if dbc := db.First(vertical, "name=?", "Shirts"); dbc.Error == nil {
	//	fmt.Println("Vertical Exists", dbc.Value.(*Vertical).Name)
	//} else {
	//	fmt.Println("Error Fetching Vertical:", dbc.Error)
	//	db.Create(&Vertical{})
	//	fmt.Println("New Vertical Created")
	//}
}
