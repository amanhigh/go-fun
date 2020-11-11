package util

import (
	"github.com/amanhigh/go-fun/apps/models/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateDbConnection(config config.Db) (db *gorm.DB, err error) {
	log.WithFields(log.Fields{"DBConfig": config}).Info("Initing DB")

	if db, err = gorm.Open(mysql.Open(config.Url), &gorm.Config{Logger: logger.Default.LogMode(config.LogLevel)}); err == nil {
		/** Print SQL */
		//db.LogMode(true)

		//TODO:From Config
		if sqlDb, err := db.DB(); err == nil {
			sqlDb.SetMaxIdleConns(config.MaxIdle)
			sqlDb.SetMaxOpenConns(config.MaxOpen)
		}
	}
	return
}

func CreateTestDb() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open("/tmp/test.db"), &gorm.Config{})
	return
}
