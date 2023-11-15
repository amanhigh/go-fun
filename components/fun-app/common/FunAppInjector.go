package common

import (
	"context"
	"fmt"
	"net/http"

	"github.com/amanhigh/go-fun/common/telemetry"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	config2 "github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/amanhigh/go-fun/models/interfaces"
	"github.com/etcinit/speedbump"
	"github.com/etcinit/speedbump/ginbump"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"gopkg.in/redis.v5"
	"gorm.io/gorm"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

const (
	NAMESPACE = "funapp"
)

type FunAppInjector struct {
	graph  inject.Graph
	config config2.FunAppConfig
}

func NewFunAppInjector(config config2.FunAppConfig) interfaces.ApplicationInjector {
	return &FunAppInjector{inject.Graph{}, config}
}

func (self *FunAppInjector) BuildApp() (app any, err error) {
	// Build App and Engin
	app = &handlers.FunServer{}
	engine := gin.New()

	/* Setup Telemetry */
	telemetry.InitLogger(self.config.Server.LogLevel)
	telemetry.InitTracerProvider(context.Background(), NAMESPACE, self.config.Tracing)

	/* Access Metrics */
	// TODO: Ingest to Prometheus and configure in helm
	//Visit http://localhost:8080/metrics
	prometheus := ginprometheus.NewPrometheus("gin_access")
	prometheus.ReqCntURLLabelMappingFn = telemetry.AccessMetrics
	prometheus.Use(engine)

	/* Middleware */
	engine.Use(gin.Recovery(), telemetry.RequestId, gin.LoggerWithFormatter(telemetry.GinRequestIdFormatter))
	// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/example/server.go
	engine.Use(otelgin.Middleware(NAMESPACE + "-gin"))

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
		&inject.Object{Value: app},
		&inject.Object{Value: &handlers.PersonHandler{}},
		&inject.Object{Value: &handlers.AdminHandler{}},
		&inject.Object{Value: util.NewGracefulShutdown()},

		&inject.Object{Value: initDb(self.config.Db)},

		&inject.Object{Value: &manager.PersonManager{}},
		&inject.Object{Value: &dao.PersonDao{}},
		&inject.Object{Value: otel.Tracer(NAMESPACE)},

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
		log.WithFields(log.Fields{"Port": self.config.Server.Port}).Info("Injection Complete")
	}
	return
}

func initDb(config config2.Db) (db *gorm.DB) {
	db = util.CreateDb(config)

	/** Gorm AutoMigrate Schema */
	db.AutoMigrate(
		&fun.Person{},
	)
	return
}
