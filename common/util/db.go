package util

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/common"
	config2 "github.com/amanhigh/go-fun/models/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateDb(dbConfig config2.Db) (db *gorm.DB) {
	var err error

	/* Create Test DB or connect to provided DB */
	if dbConfig.Url == "" {
		db, err = CreateTestDb()
	} else {
		db, err = CreateDbConnection(dbConfig)
	}

	/* Tracing */
	if err == nil {
		// https://github.com/uptrace/opentelemetry-go-extra/tree/main/otelgorm
		err = db.Use(otelgorm.NewPlugin())
	}

	/* Migrate DB */
	if err == nil && dbConfig.AutoMigrate {
		/* GoMigrate if Source is Provided */
		if dbConfig.MigrationSource != "" {
			var m *migrate.Migrate

			// Build Source and DB Url
			sourceURL := fmt.Sprintf("file://%v", dbConfig.MigrationSource)
			dbUrl := fmt.Sprintf("mysql://%v", dbConfig.Url)

			// Run Go Migrate
			if m, err = migrate.New(sourceURL, dbUrl); err == nil {
				if err = m.Up(); err == nil {
					log.WithFields(log.Fields{"Source": sourceURL, "DB": dbUrl}).Info("Migration Complete")
				} else if err == migrate.ErrNoChange {
					//Ignore No Change
					err = nil
				}
			}
		}
	}

	if err != nil {
		log.WithFields(log.Fields{"DbConfig": dbConfig, "Error": err}).Fatal("Failed To Setup DB")
	}
	return
}

func CreateDbConnection(config config2.Db) (db *gorm.DB, err error) {
	log.WithFields(log.Fields{"DBConfig": config}).Info("Initing DB")

	if db, err = gorm.Open(mysql.Open(config.Url), &gorm.Config{Logger: logger.Default.LogMode(config.LogLevel)}); err == nil {
		/** Print SQL */
		//db.LogMode(true)

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
	//BUG: Connect Log Level to Config
	db, err = gorm.Open(sqlite.Open("file:memdb1?mode=memory&cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})
	return
}

func CreateMysqlConnection(username, password, host, dbName string, port int) (db *sql.DB, err error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)
	db, err = sql.Open("mysql", url)
	return
}

// GormErrorMapper maps GORM database errors to common HTTP errors.
//
// It takes an error as a parameter and returns a common.HttpError.
func GormErrorMapper(err error) common.HttpError {
	//Doesn't Need State hence placed in Util.
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return common.ErrNotFound
		default:
			return common.NewServerError(err)
		}
	}
	return nil

}

/*
*

	Transaction Handling
*/
func Tx(c context.Context) (tx *gorm.DB) {
	if c != nil {
		//Check If Context Has Tx
		if value := c.Value(models.CONTEXT_TX); value != nil {
			//Extract and Return
			tx = value.(*gorm.DB)
		} else {
			log.Debug("Missing Transaction In Context")
		}
	} else {
		log.Debug("Nil Context Passed")
	}
	return
}
