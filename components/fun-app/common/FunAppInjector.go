package common

import (
	"fmt"
	clients2 "github.com/amanhigh/go-fun/common/clients"
	metrics2 "github.com/amanhigh/go-fun/common/metrics"
	util2 "github.com/amanhigh/go-fun/common/util"
	handlers2 "github.com/amanhigh/go-fun/components/fun-app/handlers"
	manager2 "github.com/amanhigh/go-fun/components/fun-app/manager"
	config2 "github.com/amanhigh/go-fun/models/config"
	db3 "github.com/amanhigh/go-fun/models/fun-app/db"
	interfaces2 "github.com/amanhigh/go-fun/models/interfaces"
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
	config config2.FunAppConfig
}

func NewFunAppInjector(config config2.FunAppConfig) interfaces2.ApplicationInjector {
	return &FunAppInjector{inject.Graph{}, config}
}

func (self *FunAppInjector) BuildApp() (app interface{}, err error) {
	server := &handlers2.FunServer{}
	app = server

	/* Set Logger to File */
	var file *os.File
	if file, err = os.Create(APP_LOG); err == nil {
		log.SetOutput(file)
		//Auto Log RequestId
		log.AddHook(&metrics2.ContextLogHook{})

		/* Gin Engine */
		engine := gin.New()
		if file, err = os.Create(ACCESS_LOG); err == nil {
			engine.Use(gin.LoggerWithWriter(file), gin.Recovery(), metrics2.RequestId, metrics2.MatchedPath, metrics2.AccessMetrics)

			/* Injections */
			err = self.graph.Provide(
				&inject.Object{Value: engine},
				&inject.Object{Value: &http.Server{
					Addr:    fmt.Sprintf(":%v", self.config.Server.Port),
					Handler: engine,
				}},
				&inject.Object{Value: server},
				&inject.Object{Value: &handlers2.PersonHandler{}},
				&inject.Object{Value: &handlers2.AdminHandler{}},
				&inject.Object{Value: util2.NewGracefulShutdown()},

				&inject.Object{Value: initDb(self.config.Db)},

				&inject.Object{Value: clients2.NewHttpClient(self.config.Http)},

				&inject.Object{Value: &manager2.PersonManager{}},

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

func initDb(dbConfig config2.Db) (db *gorm.DB) {
	var err error

	/* Create Test DB or connect to provided DB */
	if dbConfig.Url == "" {
		db, err = util2.CreateTestDb()
	} else {
		db, err = util2.CreateDbConnection(dbConfig)
	}

	/* Migrate DB */
	if err == nil && dbConfig.AutoMigrate {
		/** Gorm AutoMigrate Schema */
		db.AutoMigrate(
			&db3.Person{},
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
