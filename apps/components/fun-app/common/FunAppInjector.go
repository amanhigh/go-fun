package common

import (
	"fmt"
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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	APP_LOG    = "/var/log/fun-app/service.log"
	ACCESS_LOG = "/var/log/fun-app/access.log"
	NAMESPACE  = "fun_app"
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
	if file, err = os.Create(APP_LOG); err == nil {
		log.SetOutput(file)
		//Auto Log RequestId
		log.AddHook(&metrics.ContextLogHook{})

		/* Gin Engine */
		engine := gin.New()
		if file, err = os.Create(ACCESS_LOG); err == nil {
			engine.Use(gin.LoggerWithWriter(file), gin.Recovery(), metrics.RequestId, metrics.MatchedPath, metrics.AccessMetrics)

			/* Injections */
			err = self.graph.Provide(
				&inject.Object{Value: engine},
				&inject.Object{Value: &http.Server{
					Addr:    fmt.Sprintf(":%v", self.config.Server.Port),
					Handler: engine,
				}},
				&inject.Object{Value: server},
				&inject.Object{Value: &handlers.PersonHandler{}},
				&inject.Object{Value: &handlers.AdminHandler{}},
				&inject.Object{Value: util.NewGracefulShutdown()},

				&inject.Object{Value: initDb(self.config.Db)},

				&inject.Object{Value: clients.NewHttpClient(self.config.Http)},

				&inject.Object{Value: &manager.PersonManager{}},

				/* Metrics */
				&inject.Object{Value: promauto.NewCounterVec(prometheus.CounterOpts{
					Namespace:   NAMESPACE,
					Name:        "create_person",
					Help:        "Counts Person Create API",
					ConstLabels: nil,
				}, []string{"gender"}), Name: "m_create_person"},
				&inject.Object{Value: promauto.NewGauge(prometheus.GaugeOpts{
					Namespace:   NAMESPACE,
					Name:        "person_count",
					Help:        "Person Count in Get Persons",
					ConstLabels: nil,
				}), Name: "m_person_count"},
				&inject.Object{Value: promauto.NewHistogram(prometheus.HistogramOpts{
					Namespace:   NAMESPACE,
					Name:        "person_create_time",
					Help:        "Time Taken to Create Person",
					ConstLabels: nil,
				}), Name: "m_person_create_time"},
			)

			if err == nil {
				err = self.graph.Populate()
			}
		}

		//go metrics.WriteToFile("metrics.prom",10 * time.Second)
	}
	return
}

func initDb(dbConfig config.Db) (db *gorm.DB) {
	var err error

	/* Create Test DB or connect to provided DB */
	if dbConfig.Url == "" {
		db, err = util.CreateTestDb()
	} else {
		db, err = util.CreateDbConnection(dbConfig)
	}

	/* Migrate DB */
	if err == nil && dbConfig.AutoMigrate {
		/** Gorm AutoMigrate Schema */
		db.AutoMigrate(
			&db2.Person{},
		)

		/* GoMigrate*/
		if dbConfig.MigrationSource != "" {
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
