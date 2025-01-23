package common

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	"github.com/rs/zerolog/log"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	metric_sdk "go.opentelemetry.io/otel/sdk/metric"

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
	/* Setup Telemetry */
	telemetry.InitLogger(self.config.Log)
	telemetry.InitTracerProvider(context.Background(), NAMESPACE, self.config.Tracing)
	setupPrometheus()

	/* Validators */
	v, _ := binding.Validator.Engine().(*validator.Validate)
	_ = v.RegisterValidation("name", NameValidator)

	/* Injections */
	// Build Libraries
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
	registerMetrics(self.di)

	//Build Components
	container.MustSingleton(self.di, util.NewBaseDao)
	container.MustSingleton(self.di, dao.NewPersonDao)
	container.MustSingleton(self.di, manager.NewPersonManager)

	registerHandlers(self.di)

	// Build App
	app = &handlers.FunServer{}
	err = self.di.Fill(app)
	if err == nil {
		log.Info().Int("Port", self.config.Server.Port).Msg("Injection Complete")
	}
	return
}

func newHttpServer(config config.FunAppConfig, engine *gin.Engine) (server *http.Server) {
	server = &http.Server{
		Addr:              fmt.Sprintf(":%v", config.Server.Port),
		Handler:           engine,
		ReadHeaderTimeout: 10 * time.Second,
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
		log.Info().Str("Redis", cfg.RedisHost).Int64("RateLimit", cfg.PerMinuteLimit).Msg("Rate Limit Enabled")
	}
}

func setupPrometheus() {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Prometheus Exporter Failed")
	}

	provider := metric_sdk.NewMeterProvider(
		metric_sdk.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)
}

func newPrometheus(engine *gin.Engine) (prometheus *ginprometheus.Prometheus) {
	/* Access Metrics */
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
		&fun.PersonAudit{},
	)
	return
}

func registerMetrics(di container.Container) {
	meter := otel.GetMeterProvider().Meter(NAMESPACE)

	container.MustNamedSingleton(di, "CreateCounter", func() metric.Int64Counter {
		counter, _ := meter.Int64Counter("create_person",
			metric.WithDescription("Counts Person Create API"),
		)
		return counter
	})

	container.MustNamedSingleton(di, "PersonCounter", func() metric.Int64UpDownCounter {
		counter, _ := meter.Int64UpDownCounter("person_count",
			metric.WithDescription("Person Count in Get Persons"),
		)
		return counter
	})

	container.MustNamedSingleton(di, "PersonCreateTime", func() metric.Float64Histogram {
		histogram, _ := meter.Float64Histogram("person_create_time",
			metric.WithDescription("Time Taken to Create Person"),
		)
		return histogram
	})
}

func registerHandlers(di container.Container) {
	container.MustSingleton(di, handlers.NewAdminHandler)
	container.MustSingleton(di, func() (handler *handlers.PersonHandler, err error) {
		handler = &handlers.PersonHandler{}

		err = di.Fill(handler)
		return
	})
}
