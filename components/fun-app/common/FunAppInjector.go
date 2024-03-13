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
	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/amanhigh/go-fun/models/interfaces"
	"github.com/etcinit/speedbump"
	"github.com/etcinit/speedbump/ginbump"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golobby/container/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"gopkg.in/redis.v5"
	"gorm.io/gorm"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	NAMESPACE = "funapp"
)

type FunAppInjector struct {
	di     container.Container
	config config.FunAppConfig
}

func NewFunAppInjector(cfg config.FunAppConfig) interfaces.ApplicationInjector {
	return &FunAppInjector{container.New(), cfg}
}

func (self *FunAppInjector) BuildApp() (app any, err error) {
	// Build App and Engin
	app = &handlers.FunServer{}

	/* Setup Telemetry */
	telemetry.InitLogger(self.config.Server.LogLevel)
	telemetry.InitTracerProvider(context.Background(), NAMESPACE, self.config.Tracing)

	/* Validators */
	v, _ := binding.Validator.Engine().(*validator.Validate)
	_ = v.RegisterValidation("name", NameValidator)

	/* Injections */
	container.MustSingleton(self.di, func() config.FunAppConfig {
		return self.config
	})
	container.MustSingleton(self.di, newGin)
	container.MustSingleton(self.di, newPrometheus)
	container.MustSingleton(self.di, newHttpServer)
	container.MustSingleton(self.di, util.NewGracefulShutdown)
	container.MustSingleton(self.di, newDb)

	container.MustSingleton(self.di, func() trace.Tracer {
		return otel.Tracer(NAMESPACE)
	})

	container.MustSingleton(self.di, util.NewBaseDao)
	container.MustSingleton(self.di, dao.NewPersonDao)
	container.MustSingleton(self.di, handlers.NewAdminHandler)

	registerMetrics(self.di)
	registerHandlers(self.di)

	err = self.di.Fill(app)
	if err == nil {
		log.WithFields(log.Fields{"Port": self.config.Server.Port}).Info("Injection Complete")
	}
	return
}

func newHttpServer(config config.FunAppConfig, engine *gin.Engine) (server *http.Server) {
	server = &http.Server{
		Addr:    fmt.Sprintf(":%v", config.Server.Port),
		Handler: engine,
	}
	return
}

func newGin(config config.FunAppConfig) (engine *gin.Engine) {
	engine = gin.New()

	/* Middleware */
	engine.Use(gin.Recovery(), telemetry.RequestId, gin.LoggerWithFormatter(telemetry.GinRequestIdFormatter))
	// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/example/server.go
	engine.Use(otelgin.Middleware(NAMESPACE + "-gin"))

	/* Setup Rate Limit if enabled */
	setupRateLimit(config.RateLimit, engine)
	return
}

// setupRateLimit enables rate limiting if the limit is above 0.
//
// It takes in a config struct (config2.RateLimit) and a gin engine (*gin.Engine) as parameters.
// There is no return type for this function.
func setupRateLimit(cfg config.RateLimit, engine *gin.Engine) {
	/* Enable Rate Limit if Limit is above 0 */
	if cfg.PerMinuteLimit > 0 {
		// Create a Redis client
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.RedisHost,
			Password: "",
			DB:       0,
		})

		// Limit the engine's (Global) or group's (API Level) requests to
		// 100 requests per client per minute.
		engine.Use(ginbump.RateLimit(client, speedbump.PerMinuteHasher{}, cfg.PerMinuteLimit))
		log.WithFields(log.Fields{"Redis": cfg.RedisHost, "RateLimit": cfg.PerMinuteLimit}).Info("Rate Limit Enabled")
	}
}

func newPrometheus(engine *gin.Engine) (prometheus *ginprometheus.Prometheus) {
	/* Access Metrics */
	// TODO: #C Ingest to Prometheus and configure in helm
	//Visit http://localhost:8080/metrics
	prometheus = ginprometheus.NewPrometheus("gin_access")
	prometheus.ReqCntURLLabelMappingFn = telemetry.AccessMetrics
	prometheus.Use(engine)
	return
}

func newDb(config config.FunAppConfig) (db *gorm.DB, err error) {
	db = util.MustCreateDb(config.Db)

	/** Gorm AutoMigrate Schema */
	err = db.AutoMigrate(
		&fun.Person{},
	)
	return
}

func registerMetrics(di container.Container) {
	// 	/* Metrics */
	// 	//FIXME: Move to OTEL SDK
	container.MustNamedSingleton(di, "CreateCounter", func() *prometheus.CounterVec {
		return promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace:   NAMESPACE,
			Name:        "create_person",
			Help:        "Counts Person Create API",
			ConstLabels: nil,
		}, []string{"gender"})
	})

	container.MustNamedSingleton(di, "PersonCounter", func() prometheus.Gauge {
		return promauto.NewGauge(prometheus.GaugeOpts{
			Namespace:   NAMESPACE,
			Name:        "person_count",
			Help:        "Person Count in Get Persons",
			ConstLabels: nil,
		})
	})

	container.MustNamedSingleton(di, "PersonCreateTime", func() prometheus.Histogram {
		return promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace:   NAMESPACE,
			Name:        "person_create_time",
			Help:        "Time Taken to Create Person",
			ConstLabels: nil,
		})
	})
}

func registerHandlers(di container.Container) {

	container.MustSingleton(di, manager.NewPersonManager)

	container.MustSingleton(di, func() (handler *handlers.PersonHandler, err error) {
		handler = &handlers.PersonHandler{}

		err = di.Fill(handler)
		return
	})
}
