package common

import (
	"fmt"
	"net/http"

	metrics2 "github.com/amanhigh/go-fun/common/metrics"
	util2 "github.com/amanhigh/go-fun/common/util"
	handlers2 "github.com/amanhigh/go-fun/components/fun-app/handlers"
	manager2 "github.com/amanhigh/go-fun/components/fun-app/manager"
	config2 "github.com/amanhigh/go-fun/models/config"
	db3 "github.com/amanhigh/go-fun/models/fun-app/db"
	interfaces2 "github.com/amanhigh/go-fun/models/interfaces"
	"github.com/etcinit/speedbump"
	"github.com/etcinit/speedbump/ginbump"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"gopkg.in/redis.v5"
	"gorm.io/gorm"
)

const (
	NAMESPACE = "fun_app"
)

type FunAppInjector struct {
	graph  inject.Graph
	config config2.FunAppConfig
}

func NewFunAppInjector(config config2.FunAppConfig) interfaces2.ApplicationInjector {
	return &FunAppInjector{inject.Graph{}, config}
}

func (self *FunAppInjector) BuildApp() (app any, err error) {
	server := &handlers2.FunServer{}
	app = server

	//Auto Log RequestId
	log.AddHook(&metrics2.ContextLogHook{})
	log.SetLevel(self.config.Server.LogLevel)

	/* Gin Engine */
	engine := gin.New()

	/* Access Metrics */
	// TODO: Ingest to Prometheus and configure in helm
	//Visit http://localhost:8080/metrics
	prometheus := ginprometheus.NewPrometheus("gin_access")
	prometheus.ReqCntURLLabelMappingFn = metrics2.AccessMetrics
	prometheus.Use(engine)

	// http://localhost:8080/debug/statsviz/
	engine.GET("/debug/statsviz/*filepath", metrics2.StatvizMetrics)

	/* Middleware */
	engine.Use(gin.Recovery(), metrics2.RequestId, gin.LoggerWithFormatter(metrics2.GinRequestIdFormatter))

	/* Validators */
	v, _ := binding.Validator.Engine().(*validator.Validate)
	_ = v.RegisterValidation("name", NameValidator)

	/* Enable Rate Limit if Limit is above 0 */
	if self.config.RateLimit.PerMinuteLimit > 0 {
		// Create a Redis client
		client := redis.NewClient(&redis.Options{
			Addr:     self.config.RateLimit.RedisHost,
			Password: "",
			DB:       0,
		})

		// Limit the engine's (Global) or group's (API Level) requests to
		// 100 requests per client per minute.
		engine.Use(ginbump.RateLimit(client, speedbump.PerMinuteHasher{}, self.config.RateLimit.PerMinuteLimit))
		log.WithFields(log.Fields{"Redis": self.config.RateLimit.RedisHost, "RateLimit": self.config.RateLimit.PerMinuteLimit}).Info("Rate Limit Enabled")
	}

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

		&inject.Object{Value: &manager2.PersonManager{}},

		/* Metrics */
		&inject.Object{Value: promauto.NewCounterVec(prometheus2.CounterOpts{
			Namespace:   NAMESPACE,
			Name:        "create_person",
			Help:        "Counts Person Create API",
			ConstLabels: nil,
		}, []string{"gender"}), Name: "m_create_person"},
		&inject.Object{Value: promauto.NewGauge(prometheus2.GaugeOpts{
			Namespace:   NAMESPACE,
			Name:        "person_count",
			Help:        "Person Count in Get Persons",
			ConstLabels: nil,
		}), Name: "m_person_count"},
		&inject.Object{Value: promauto.NewHistogram(prometheus2.HistogramOpts{
			Namespace:   NAMESPACE,
			Name:        "person_create_time",
			Help:        "Time Taken to Create Person",
			ConstLabels: nil,
		}), Name: "m_person_create_time"},
	)
	if err == nil {
		err = self.graph.Populate()
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
