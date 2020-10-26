package fun_app

import (
	"fmt"
	"gorm.io/gorm"
	"os"

	"github.com/amanhigh/go-fun/apps/common/clients"
	"github.com/amanhigh/go-fun/apps/common/metrics"
	"github.com/amanhigh/go-fun/apps/common/util"
	"github.com/amanhigh/go-fun/apps/components/fun-app/handlers"
	"github.com/amanhigh/go-fun/apps/components/fun-app/manager"
	"github.com/amanhigh/go-fun/apps/models/config"
	db2 "github.com/amanhigh/go-fun/apps/models/fun-app/db"
	"github.com/amanhigh/go-fun/apps/models/interfaces"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	APP_LOG    = "/var/log/fun-app/service.log"
	ACCESS_LOG = "/var/log/fun-app/access.log"
)

type FunAppInjector struct {
	graph  inject.Graph
	config config.FunAppConfig
}

func NewFunAppInjector(config config.FunAppConfig) interfaces.ApplicationInjector {
	return &FunAppInjector{inject.Graph{}, config}
}

func (self *FunAppInjector) BuildApp() (app interface{}, err error) {
	server := &handlers.FunServer{}
	app = server

	/* Set Logger to File */
	var file *os.File
	if file, err = createLogfile(APP_LOG); err == nil {
		log.SetOutput(file)

		/* Gin Engine */
		engine := gin.New()
		if file, err = createLogfile(ACCESS_LOG); err == nil {
			engine.Use(gin.LoggerWithWriter(file), gin.Recovery(), metrics.RequestId, metrics.MatchedPath, metrics.AccessMetrics)

			/* Injections */
			err = self.graph.Provide(
				&inject.Object{Value: engine},
				&inject.Object{Value: self.config.Server.Port, Name: "port"},
				&inject.Object{Value: server},
				&inject.Object{Value: &handlers.PersonHandler{}},

				&inject.Object{Value: initDb(self.config.Db)},

				&inject.Object{Value: clients.NewHttpClient(self.config.Http)},

				&inject.Object{Value: &manager.PersonManager{}},
			)

			if err == nil {
				err = self.graph.Populate()
			}
		}
	}
	return
}

func createLogfile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

func initDb(dbConfig config.Db) (db *gorm.DB) {
	var err error
	if db, err = util.CreateDbConnection(dbConfig); err == nil {
		if dbConfig.AutoMigrate {
			/** Gorm AutoMigrate Schema */
			db.AutoMigrate(
				&db2.Person{},
			)

			/* GoMigrate*/
			var m *migrate.Migrate
			sourceURL := fmt.Sprintf("file://%v", dbConfig.MigrationSource)
			dbUrl := fmt.Sprintf("mysql://%v", dbConfig.Url)
			if m, err = migrate.New(sourceURL, dbUrl); err == nil {
				if err = m.Up(); err == nil {
					log.Info("Migration Complete")
				} else if err == migrate.ErrNoChange {
					//Ignore No Change
					err = nil
				}
			}
		}
	}

	if err != nil {
		log.WithFields(log.Fields{"DbConfig": dbConfig, "Error": err}).Panic("Failed To Setup DB")
	}
	return
}
