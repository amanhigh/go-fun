package util

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const (
	DB_TYPE = "mysql"
)

func CreateDbConnection(env, path string) (db *gorm.DB, err error) {
	log.WithFields(log.Fields{"Path": path, "Env": env}).Info("Initing DB")
	var dbConf *goose.DBConf
	if dbConf, err = goose.NewDBConf(path, env, DB_TYPE); err == nil {
		if db, err = gorm.Open(dbConf.PgSchema, dbConf.Driver.OpenStr); err == nil {
			/** Print SQL */
			//db.LogMode(true)

			//TODO:From Config
			db.DB().SetMaxIdleConns(5)
			db.DB().SetMaxOpenConns(20)
		}
	}
	return
}
