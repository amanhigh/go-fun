package util

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/glebarez/sqlite"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateDb(cfg config.Db) (db *gorm.DB, err error) {
	/* Create Test DB or connect to provided DB */
	if cfg.Url == "" {
		db, err = CreateTestDb()
	} else {
		db, err = ConnectDb(cfg)
	}

	/* Tracing */
	if err == nil {
		// https://github.com/uptrace/opentelemetry-go-extra/tree/main/otelgorm
		err = db.Use(otelgorm.NewPlugin())
	}

	/* Migrate DB */
	if err == nil && cfg.AutoMigrate {
		/* GoMigrate if Source is Provided */
		if cfg.MigrationSource != "" {
			var m *migrate.Migrate

			// Build Source and DB Url
			sourceURL := fmt.Sprintf("file://%v", cfg.MigrationSource)
			dbUrl := fmt.Sprintf("mysql://%v", cfg.Url)

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
	return
}

func MustCreateDb(cfg config.Db) *gorm.DB {
	db, err := CreateDb(cfg)
	if err != nil {
		log.WithFields(log.Fields{"DbConfig": cfg, "Error": err}).Fatal("Failed To Setup DB")
	}
	return db
}

func ConnectDb(cfg config.Db) (db *gorm.DB, err error) {
	log.WithFields(log.Fields{"DBConfig": cfg}).Info("Initing DB")

	var dialector gorm.Dialector

	switch strings.ToLower(cfg.DbType) {
	case "postgres":
		dialector = postgres.Open(cfg.Url)
	case "mysql":
		dialector = mysql.Open(cfg.Url)
	default:
		return nil, fmt.Errorf("unsupported db type: %s", cfg.DbType)
	}

	// FIXME: #C Support Postgress
	if db, err = gorm.Open(dialector, &gorm.Config{Logger: logger.Default.LogMode(cfg.LogLevel)}); err == nil {
		/** Print SQL */
		//db.LogMode(true)

		if sqlDb, err := db.DB(); err == nil {
			sqlDb.SetMaxIdleConns(cfg.MaxIdle)
			sqlDb.SetMaxOpenConns(cfg.MaxOpen)
		}
	}
	return
}

// CreateTestDb creates a test database.
// Uses in memory Sqlite.
// Faster CGO Implementation - https://github.com/go-gorm/sqlite
// It returns a *gorm.DB and an error.
func CreateTestDb() (db *gorm.DB, err error) {
	//Use Log Level 4 for Debug, 3 for Warnings, 2 for Errors
	//Can use /tmp/gorm.db for file base Db
	//BUG: Connect Log Level to Config
	db, err = gorm.Open(sqlite.Open("file:memdb1?mode=memory&cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})
	return
}

func CreateMysqlConnection(username, password, host, dbName string, port int) (db *sql.DB, err error) {
	url := BuildMysqlURL(username, password, host, dbName, port)
	db, err = sql.Open("mysql", url)
	return
}

func BuildMysqlURL(username, password, host, dbName string, port any) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)
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
