package core

import (
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
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

// registerBarkatDependencies registers all dependencies for the Barkat ticker feature.
func (ki *KohanInjector) registerBarkatDependencies() error {
	// Ticker
	container.MustSingleton(ki.di, provideTickerRepository)
	container.MustSingleton(ki.di, provideBarkatTickerManager)
	container.MustSingleton(ki.di, provideTickerHandler)

	return nil
}
