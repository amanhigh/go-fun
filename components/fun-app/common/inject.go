package common

import (
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/interfaces"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golobby/container/v3"
	"github.com/rs/zerolog/log"
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

func (fi *FunAppInjector) BuildApp() (app any, err error) {
	fi.setupTelemetry()
	fi.registerValidators()

	// Register dependencies bottom-up: infra → data → domain → handlers → messaging
	fi.registerCoreDependencies()
	fi.registerMetrics()
	fi.registerMessagingInfra()
	fi.registerDao()
	fi.registerPublishers()
	fi.registerManager()
	fi.registerHandlers()
	fi.registerCommandHandlers()
	fi.registerMessagingWiring()

	app, err = fi.buildApplication()
	return
}

func (fi *FunAppInjector) registerCoreDependencies() {
	container.MustSingleton(fi.di, func() config.FunAppConfig {
		return fi.config
	})

	// HACK: Don't inject Rate & Http Config
	container.MustSingleton(fi.di, func() config.RateLimit {
		return fi.config.RateLimit
	})
	container.MustSingleton(fi.di, func() config.HttpServerConfig {
		return config.HttpServerConfig{Name: NAMESPACE, Port: fi.config.Server.Port}
	})
	container.MustSingleton(fi.di, util.NewGracefulShutdown)
	container.MustSingleton(fi.di, newBaseHTTPServer)
	container.MustSingleton(fi.di, newPrometheus)
	container.MustSingleton(fi.di, newDb)
	container.MustSingleton(fi.di, func() trace.Tracer {
		return otel.Tracer(NAMESPACE)
	})
}

func (fi *FunAppInjector) registerValidators() {
	v, _ := binding.Validator.Engine().(*validator.Validate)
	_ = v.RegisterValidation("name", NameValidator)
}

func (fi *FunAppInjector) buildApplication() (app any, err error) {
	// HACK: Streamline Base Server Building.
	var base *util.HttpServer
	if err = fi.di.Resolve(&base); err != nil {
		return nil, fmt.Errorf("failed to resolve base http server: %w", err)
	}

	lifecycle := &handlers.FunAppServerLifecycle{}
	if err = fi.di.Fill(lifecycle); err != nil {
		return nil, fmt.Errorf("failed to fill fun app lifecycle: %w", err)
	}
	base.SetLifecycle(lifecycle)

	log.Info().Int("Port", fi.config.Server.Port).Msg("Injection Complete")
	return base, nil
}
