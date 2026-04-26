package core

import (
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/config"

	"github.com/golobby/container/v3"
)

// =============================================================================
// INTERFACE DEFINITION
// =============================================================================

// KohanInterface defines the public API for the Kohan dependency injection system
type KohanInterface interface {
	GetDariusApp(cfg config.DariusConfig) (*DariusV1, error)
	GetAutoManager(wait time.Duration) manager.AutoManagerInterface
	GetTaxManager() (manager.TaxManager, error)
	GetBrokerageManager() (manager.BrokerageManager, error)
	GetKohanServer(port int, wait time.Duration) (util.HttpServer, error)
}

// =============================================================================
// INJECTOR SETUP
// =============================================================================

// Private singleton instance
var globalInjector *KohanInjector

type KohanInjector struct {
	di     container.Container
	config config.KohanConfig
}

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

// =============================================================================
// INTERFACE METHODS
// =============================================================================
// These methods implement the KohanInterface and resolve dependencies from the injector

func (ki *KohanInjector) GetAutoManager(wait time.Duration) manager.AutoManagerInterface {
	return manager.NewAutoManager(wait, ki.config.Barkat.ScreenshotPath)
}

func (ki *KohanInjector) GetKohanServer(port int, wait time.Duration) (util.HttpServer, error) {
	autoManager := ki.GetAutoManager(wait)

	// Register all dependencies
	ki.registerOSDependencies(autoManager)
	ki.registerPortalDependencies()
	if err := ki.registerJournalDependencies(); err != nil {
		return nil, fmt.Errorf("failed to register journal dependencies: %w", err)
	}
	ki.registerServerDependencies(port)

	// Resolve server from DI
	var server util.HttpServer
	if err := ki.di.Resolve(&server); err != nil {
		return nil, fmt.Errorf("failed to resolve server: %w", err)
	}

	return server, nil
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

// =============================================================================
// PROVIDER METHODS
// =============================================================================
// These methods register dependencies with the container

func (ki *KohanInjector) registerServerDependencies(port int) {
	container.MustSingleton(ki.di, util.NewGracefulShutdown)
	container.MustSingleton(ki.di, func() config.HttpServerConfig {
		return config.HttpServerConfig{Name: "kohan", Port: port}
	})
	container.MustSingleton(ki.di, provideHttpServer)
}
