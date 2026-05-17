package core

import (
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/golobby/container/v3"
	"gorm.io/gorm"
)

// ---- Ticker Providers ----

func provideTickerRepository(db *gorm.DB) repository.TickerRepository {
	return repository.NewTickerRepository(db)
}

func provideBarkatTickerManager(repo repository.TickerRepository) manager.BarkatTickerManager {
	return manager.NewBarkatTickerManager(repo)
}

func provideTickerHandler(mgr manager.BarkatTickerManager) handler.TickerHandler {
	return handler.NewTickerHandler(mgr)
}

// ---- Alert Ticker Providers ----

func provideAlertTickerRepository(db *gorm.DB) repository.AlertTickerRepository {
	return repository.NewAlertTickerRepository(db)
}

func provideAlertTickerManager(repo repository.AlertTickerRepository) manager.AlertTickerManager {
	return manager.NewAlertTickerManager(repo)
}

func provideAlertTickerHandler(mgr manager.AlertTickerManager) handler.AlertTickerHandler {
	return handler.NewAlertTickerHandler(mgr)
}

// registerBarkatDependencies registers all dependencies for the Barkat ticker feature.
func (ki *KohanInjector) registerBarkatDependencies() error {
	// Database — must be registered first since all repos depend on it.
	container.MustSingleton(ki.di, func() config.BarkatConfig { return ki.config.Barkat })

	// Eagerly create DB so we get a clear error if it fails, then register the concrete instance.
	db, err := ki.provideBarkatDB()
	if err != nil {
		return fmt.Errorf("failed to open barkat database: %w", err)
	}
	container.MustSingleton(ki.di, func() *gorm.DB { return db })

	// Ticker
	container.MustSingleton(ki.di, provideTickerRepository)
	container.MustSingleton(ki.di, provideBarkatTickerManager)
	container.MustSingleton(ki.di, provideTickerHandler)

	// Alert Ticker
	container.MustSingleton(ki.di, provideAlertTickerRepository)
	container.MustSingleton(ki.di, provideAlertTickerManager)
	container.MustSingleton(ki.di, provideAlertTickerHandler)

	return nil
}
