package util

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func NewDb(path string) (db *gorm.DB) {
	var err error
	var dbConf *goose.DBConf
	if dbConf, err = goose.NewDBConf(path, "development", "mysql"); err == nil {
		db, err = gorm.Open(dbConf.PgSchema, dbConf.Driver.OpenStr)
		//db.LogMode(true)

		db.DB().SetMaxIdleConns(5)
		db.DB().SetMaxOpenConns(20)

		//goose.RunMigrations(dbConf,path,1)
	}

	if err != nil {
		log.WithFields(log.Fields{"Path": path, "User": "root", "Type": "mysql"}).Panic("failed to connect database")
	}
	return
}
