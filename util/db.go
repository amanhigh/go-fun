package util

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

func NewDb(path string) (db *gorm.DB) {
	var err error
	var dbConf *goose.DBConf
	if dbConf, err = goose.NewDBConf(path, "development", "mysql"); err == nil {
		db, err = gorm.Open(dbConf.PgSchema, dbConf.Driver.OpenStr)
		//db.LogMode(true)
	}

	if err != nil {
		log.WithFields(log.Fields{"Path": path, "User": "root", "Type": "mysql"}).Panic("failed to connect database")
	}
	return
}
