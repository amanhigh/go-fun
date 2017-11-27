package util

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

func NewDb(path string) *gorm.DB {
	dbConf, _ := goose.NewDBConf(path, "development", "mysql")

	if db, err := gorm.Open(dbConf.PgSchema, dbConf.Driver.OpenStr); err != nil {
		log.WithFields(log.Fields{"Path": path, "User": "root", "Type": "mysql"}).Panic("failed to connect database")
		return nil
	} else {
		/** Print SQL */
		if IsDebugMode() {
			db.LogMode(true)
		}
		return db
	}
}

func Migrate(db *gorm.DB, values ...interface{}) {
	/** AutoMigrate Schema */
	db.AutoMigrate(values)
}
