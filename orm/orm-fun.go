package main

import "github.com/jinzhu/gorm"
import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
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

func main() {
	db, err := gorm.Open("mysql", "root@/aman?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	migrate(db)

	playProduct(db)

	//db.Create(&Product{Code: "LongCode", Price: 4})

	//schemaAlterPlay(db)
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

	// find product with code l1212
	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)
	// Delete - delete product
	db.Delete(&product)
}
func queryProduct(db *gorm.DB) {
	// First Query
	product := new(Product)
	db.First(product, 1)

	// Preload with Where Clause
	db.Preload("Vertical").First(product, "code = ?", "L1212")
	fmt.Println("Vertical ID:", product.VerticalId, "Vertical Name:", product.Vertical.Name)

	//Query all Non Deleted Products
	products := new([]Product)
	db.Unscoped().Where("code = ?", "L1212").Find(products)
	for _, product := range *products {
		fmt.Println("Deleted/Undeleted Product Found: ", product.ID)
	}

	//Query Id Range
	db.Unscoped().Not([]int64{5, 6, 10}).Find(products)
	db.Unscoped().Where([]int64{5, 6, 10}).Find(products)
	fmt.Println("Id Range Search Count: ", len(*products))

	//Struct Query
	db.Where(&Product{Price: 2000}).Last(product)
	fmt.Println("Query By Struct, ID:", product.ID)
}

func migrate(db *gorm.DB) {
	/** Clear Old Tables */
	//db.DropTable(&Product{}, &Vertical{})
	// Migrate the schema
	db.AutoMigrate(&Product{}, &Vertical{})
}

func createVertical(db *gorm.DB) {
	vertical := &Vertical{}
	if dbc := db.First(vertical, "name=?", "Shirts"); dbc.Error == nil {
		fmt.Println("Vertical Exists", dbc.Value.(*Vertical).Name)
	} else {
		fmt.Println("Error Fetching Vertical:", dbc.Error)
		db.Create(&Vertical{})
		fmt.Println("New Vertical Created")
	}
}
