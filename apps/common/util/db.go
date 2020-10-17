package util

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateDbConnection(env, path string) (db *gorm.DB, err error) {
	log.WithFields(log.Fields{"Path": path, "Env": env}).Info("Initing DB")

	if db, err = gorm.Open(mysql.Open(path), &gorm.Config{}); err == nil {
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
