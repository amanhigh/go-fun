package util

import (
	"database/sql"
	"fmt"
	"github.com/amanhigh/go-fun/apps/models/config"
	_ "github.com/go-sql-driver/mysql"
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
	//Use Log Level 4 for Debug, 3 for Warnings, 2 for Errors
	//Can use /tmp/gorm.db for file base Db
	db, err = gorm.Open(sqlite.Open("file:memdb1?mode=memory&cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})
	return
}

func CreateMysqlConnection(username, password, host, dbName string, port int) (db *sql.DB, err error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)
	db, err = sql.Open("mysql", url)
	return
}
