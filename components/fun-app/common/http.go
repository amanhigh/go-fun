package common

import (
	"github.com/amanhigh/go-fun/common/telemetry"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/etcinit/speedbump"
	"github.com/etcinit/speedbump/ginbump"
	"github.com/gin-gonic/gin"
	"github.com/golobby/container/v3"
	"github.com/rs/zerolog/log"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gopkg.in/redis.v5"
)

func newGin(rateCfg config.RateLimit) (engine *gin.Engine) {
	engine = gin.New()

	/* Middleware */
	engine.Use(gin.Recovery(), telemetry.RequestID, gin.LoggerWithFormatter(telemetry.GinRequestIdFormatter))
	// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/example/server.go
	engine.Use(otelgin.Middleware(NAMESPACE + "-gin"))

	/* Setup Rate Limit if enabled */
	setupRateLimit(rateCfg, engine)
	return
}

func provideHttpServer(cfg config.HttpServerConfig, rateCfg config.RateLimit, shutdown util.Shutdown) *util.HttpServer {
	engine := newGin(rateCfg)
	return util.NewHttpServer(cfg, engine, shutdown)
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

func newPrometheus(base *util.HttpServer) (prometheus *ginprometheus.Prometheus) {
	// HACK: Only pass required engine here.
	/* Access Metrics */
	// Visit http://localhost:8080/metrics
	prometheus = ginprometheus.NewPrometheus("gin_access")
	prometheus.ReqCntURLLabelMappingFn = telemetry.AccessMetrics
	prometheus.Use(base.Engine)
	return
}

func (fi *FunAppInjector) registerHandlers() {
	container.MustSingleton(fi.di, handlers.NewAdminHandler)
	container.MustSingleton(fi.di, fi.providePersonHandler)
	container.MustSingleton(fi.di, handlers.NewEnrollmentHandler)
}

func (fi *FunAppInjector) providePersonHandler() (handler handlers.PersonHandler, err error) {
	impl := &handlers.PersonHandlerImpl{}
	err = fi.di.Fill(impl)
	handler = impl
	return
}
