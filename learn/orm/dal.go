package orm

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

	/** Print SQL */
	//db.LogMode(true)

	return db
}

func Migrate(db *gorm.DB,values ...interface{}) {
	/** Clear Old Tables */
	//db.DropTable(&Product{}, &Vertical{})

	/** AutoMigrate Schema */
	db.AutoMigrate(values)
}