package core

import (
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/config"

	"github.com/golobby/container/v3"
	"gorm.io/gorm"
)

// Interface and implementation in same file
type KohanInterface interface {
	GetDariusApp(cfg config.DariusConfig) (*DariusV1, error)
	GetAutoManager(wait time.Duration, capturePath string) manager.AutoManagerInterface
	GetTaxManager() (manager.TaxManager, error)
	GetBrokerageManager() (manager.BrokerageManager, error)
	GetKohanServer(port int, capturePath string, wait time.Duration, shutdown util.Shutdown) (*KohanServer, error)
	GetBarkatDB() (*gorm.DB, error)
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

func (ki *KohanInjector) GetKohanServer(port int, capturePath string, wait time.Duration, shutdown util.Shutdown) (*KohanServer, error) {
	// HACK: Shutdown should be created internally not passed as parameter.
	autoManager := ki.GetAutoManager(wait, capturePath)
	ki.registerMonitorDependencies(capturePath, autoManager)
	if err := ki.registerJournalDependencies(); err != nil {
		return nil, fmt.Errorf("failed to register journal dependencies: %w", err)
	}

	base := provideBaseHTTPServer(port, shutdown)

	server := &KohanServer{BaseHTTPServer: base}
	if err := ki.di.Fill(server); err != nil {
		return nil, fmt.Errorf("failed to fill kohan server: %w", err)
	}

	server.RegisterRoutes = server.registerRoutes
	return server, nil
}

func (ki *KohanInjector) GetBarkatDB() (*gorm.DB, error) {
	// HACK: Why we need Public method for this remove it.
	// FIXME: DB Migration has many indexes on Primary key remove unwanted indexes.
	var db *gorm.DB
	if err := ki.di.Resolve(&db); err != nil {
		return nil, fmt.Errorf("failed to resolve barkat db: %w", err)
	}
	return db, nil
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
