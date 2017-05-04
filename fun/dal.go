package fun

import (
	"github.com/jinzhu/gorm"
	log "github.com/Sirupsen/logrus"
	"bitbucket.org/liamstask/goose/lib/goose"
)

var DB=build()

func build() (*gorm.DB) {
	dbConf, _ := goose.NewDBConf("/Users/amanpreet.singh/IdeaProjects/GoArena/src/github.com/amanhigh/go-fun/orm/db/", "development", "mysql")
	db, err := gorm.Open(dbConf.PgSchema, dbConf.Driver.OpenStr)
	if err != nil {
		log.WithFields(log.Fields{
			"DB":   "aman",
			"User": "root",
			"Type": "mysql",
		}).Panic("failed to connect database")
	}
	return db
}

func Migrate(db *gorm.DB,values ...interface{}) {
	/** Print SQL */
	//db.LogMode(true)
	/** Clear Old Tables */
	//db.DropTable(&Product{}, &Vertical{})
	// Migrate the schema
	db.AutoMigrate(values)
}