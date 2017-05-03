package main

import "github.com/jinzhu/gorm"
import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"os"
)

type Product struct {
	gorm.Model
	Code       string `gorm:"size 5"`
	Price      uint
	IgnoreMe   string `gorm:"-"` // Ignore this field
	Vertical   Vertical
	VerticalId uint //Must be vertical_id in DB or won't work automatically.
}

type Vertical struct {
	gorm.Model
	Name string `gorm:"unique;default:'Shirts'"`
}

//Default Name would be products
func (p *Product) TableName() string {
	return "MeraProduct"
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

func (p *Product) AfterFind() (err error) {
	p.IgnoreMe = "Ignore" + p.Code
	return nil
}

func main() {
	db, err := gorm.Open("mysql", "root@/aman?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.WithFields(log.Fields{
			"DB":   "aman",
			"User": "root",
			"Type": "mysql",
		}).Panic("failed to connect database")
	}
	defer db.Close()

	prepLogger()
	migrate(db)


	playProduct(db)

	//db.Create(&Product{Code: "LongCode", Price: 4})

	//schemaAlterPlay(db)
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
	db.Model(&Product{}).DropColumn("code")
}

func playProduct(db *gorm.DB) {
	createVertical(db)

	// Create
	product := &Product{Code: "L1212", Price: 1000, VerticalId: 1}
	db.Create(product)

	queryProduct(db)

	productUpdates(db, product)

	// Delete - delete product
	db.Delete(&product)
}
func productUpdates(db *gorm.DB, product *Product) {
	// Update without Callbacks
	db.Model(&product).UpdateColumn("code", "No Callback")
	//Single Field Update
	db.Model(&product).Update("Price", 1500)
	//Struct Update
	db.Model(&product).Update(&Product{Code: "MyCode"})
}
func queryProduct(db *gorm.DB) {
	// First Query
	product := new(Product)
	db.First(product, 1)

	// Preload with Where Clause
	db.Preload("Vertical").First(product, "code = ?", "L1212")

	logContext := log.WithFields(log.Fields{
		"Vertical ID:":  product.VerticalId,
		"Vertical Name": product.Vertical.Name,
		"Ignore Me":     product.IgnoreMe,
	})
	logContext.Info("Product Details")

	//Query all Non Deleted Products
	products := new([]Product)
	db.Unscoped().Where("code = ?", "L1212").Find(products)
	for _, product := range *products {
		fmt.Println("Deleted/Undeleted Product Found: ", product.ID)
	}

	//Query Id Range
	codes := new([]string)
	db.Not([]int64{5, 6, 10}).Find(products).Pluck("code", codes)
	fmt.Printf("CODES: %+v\n", *codes)
	db.Unscoped().Where([]int64{5, 6, 10}).Limit(3).Limit(-1).Find(products)
	fmt.Println("Id Range Search Count: ", len(*products))

	//Struct Query
	db.Order("code desc,price asc").Where(&Product{Price: 2000}).Where(&Product{Code: "L1212"}).Last(product) //And
	db.Where(&Product{Price: 2000}).Or(&Product{Code: "L1212"}).Last(product)                                 //Or
	fmt.Println("Query By Struct, ID:", product.ID)
}

func migrate(db *gorm.DB) {
	/** Print SQL */
	//db.LogMode(true)
	/** Clear Old Tables */
	//db.DropTable(&Product{}, &Vertical{})
	// Migrate the schema
	db.AutoMigrate(&Product{}, &Vertical{})
}

func createVertical(db *gorm.DB) {
	vertical := &Vertical{}
	db.FirstOrCreate(&vertical)
	verticalCount := new(int)
	db.Model(&Vertical{}).Count(verticalCount)
	fmt.Println("Vertical Count:", *verticalCount)

	found := db.Model(&Vertical{Name: "Not Present"}).RecordNotFound()
	fmt.Println("FOUND Value:", found)

	//if dbc := db.First(vertical, "name=?", "Shirts"); dbc.Error == nil {
	//	fmt.Println("Vertical Exists", dbc.Value.(*Vertical).Name)
	//} else {
	//	fmt.Println("Error Fetching Vertical:", dbc.Error)
	//	db.Create(&Vertical{})
	//	fmt.Println("New Vertical Created")
	//}
}
