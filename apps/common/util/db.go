package util

import (
	"github.com/amanhigh/go-fun/apps/models/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateDbConnection(config config.Db) (db *gorm.DB, err error) {
	log.WithFields(log.Fields{"DBConfig": config}).Info("Initing DB")

	if db, err = gorm.Open(mysql.Open(config.Url), &gorm.Config{}); err == nil {
		/** Print SQL */
		//db.LogMode(true)

		//TODO:From Config
		if sqlDb, err := db.DB(); err == nil {
			sqlDb.SetMaxIdleConns(5)
			sqlDb.SetMaxOpenConns(20)
		}
	}
	return
}
