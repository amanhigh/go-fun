package core

import (
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/config"

	"github.com/golobby/container/v3"
)

// Interface and implementation in same file
type KohanInterface interface {
	GetDariusApp(cfg config.DariusConfig) (*DariusV1, error)
	GetAutoManager(wait time.Duration, capturePath string) manager.AutoManagerInterface
	GetTaxManager() (manager.TaxManager, error)
	GetBrokerageManager() (manager.BrokerageManager, error)
	GetKohanServer(port int, capturePath string, wait time.Duration) (*util.BaseHTTPServer, error)
}

// Private singleton instance
var globalInjector *KohanInjector

type KohanInjector struct {
	di     container.Container
	config config.KohanConfig
}

// ---- KohanInjector Methods ----

func SetupKohanInjector(config config.KohanConfig) {
	globalInjector = &KohanInjector{
		di:     container.New(),
		config: config,
	}
}

// Public singleton access - returns interface only
func GetKohanInterface() KohanInterface {
	return globalInjector
}

func (ki *KohanInjector) GetAutoManager(wait time.Duration, capturePath string) manager.AutoManagerInterface {
	return manager.NewAutoManager(wait, capturePath)
}

func (ki *KohanInjector) GetKohanServer(port int, capturePath string, wait time.Duration) (*util.BaseHTTPServer, error) {
	autoManager := ki.GetAutoManager(wait, capturePath)
	ki.registerMonitorDependencies(capturePath, autoManager)
	if err := ki.registerJournalDependencies(); err != nil {
		return nil, fmt.Errorf("failed to register journal dependencies: %w", err)
	}
	ki.registerServerDependencies(port)
	// FIXME: DB Migration has many indexes on Primary key remove unwanted indexes.

	var base *util.BaseHTTPServer
	if err := ki.di.Resolve(&base); err != nil {
		return nil, fmt.Errorf("failed to resolve base http server: %w", err)
	}

	lifecycle := &KohanServerLifecycle{}
	if err := ki.di.Fill(lifecycle); err != nil {
		return nil, fmt.Errorf("failed to fill kohan lifecycle: %w", err)
	}
	base.SetLifecycle(lifecycle)
	return base, nil
}

func (ki *KohanInjector) registerServerDependencies(port int) {
	// FIXME: Sort this mess cleanly build base server and lifecycle.
	container.MustSingleton(ki.di, util.NewGracefulShutdown)
	container.MustSingleton(ki.di, func(shutdown util.Shutdown) util.HttpServerConfig {
		return util.HttpServerConfig{Name: "kohan", Port: port, Shutdown: shutdown}
	})
	container.MustSingleton(ki.di, provideBaseHTTPServer)
}

func (ki *KohanInjector) GetTaxManager() (manager.TaxManager, error) {
	ki.registerTaxDependencies()

	// Resolve and return TaxManager
	var taxManager manager.TaxManager
	err := ki.di.Resolve(&taxManager)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tax manager: %w", err)
	}
	return taxManager, nil
}

func (ki *KohanInjector) GetBrokerageManager() (manager.BrokerageManager, error) {
	ki.registerTaxDependencies()

	var brokerageManager manager.BrokerageManager
	err := ki.di.Resolve(&brokerageManager)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve brokerage manager: %w", err)
	}
	return brokerageManager, nil
}

func (ki *KohanInjector) GetDariusApp(cfg config.DariusConfig) (*DariusV1, error) {
	ki.registerDariusDependencies(cfg)

	// Build app
	app := &DariusV1{}
	err := ki.di.Fill(app)
	if err != nil {
		return nil, fmt.Errorf("failed to fill darius app: %w", err)
	}
	return app, nil
}
