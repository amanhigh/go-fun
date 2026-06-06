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

// ---- Price Alert Providers ----

func providePriceAlertRepository(db *gorm.DB) repository.PriceAlertRepository {
	return repository.NewPriceAlertRepository(db)
}

func providePriceAlertManager(repo repository.PriceAlertRepository) manager.PriceAlertManager {
	return manager.NewPriceAlertManager(repo)
}

func providePriceAlertHandler(mgr manager.PriceAlertManager) handler.PriceAlertHandler {
	return handler.NewPriceAlertHandler(mgr)
}

// ---- Audit Providers ----

func provideAuditRepository(db *gorm.DB) repository.AuditRepository {
	return repository.NewAuditRepository(db)
}

func provideAuditManager(repo repository.AuditRepository) manager.AuditManager {
	return manager.NewAuditManager(repo)
}

func provideAuditHandler(mgr manager.AuditManager) handler.AuditHandler {
	return handler.NewAuditHandler(mgr)
}

// registerBarkatDependencies registers all dependencies for the Barkat ticker feature.
func (ki *KohanInjector) registerBarkatDependencies() error {
	// Ticker
	container.MustSingleton(ki.di, provideTickerRepository)
	container.MustSingleton(ki.di, provideBarkatTickerManager)
	container.MustSingleton(ki.di, provideTickerHandler)

	// Alert Ticker
	container.MustSingleton(ki.di, provideAlertTickerRepository)
	container.MustSingleton(ki.di, provideAlertTickerManager)
	container.MustSingleton(ki.di, provideAlertTickerHandler)

	// Price Alert
	container.MustSingleton(ki.di, providePriceAlertRepository)
	container.MustSingleton(ki.di, providePriceAlertManager)
	container.MustSingleton(ki.di, providePriceAlertHandler)

	// Audit
	container.MustSingleton(ki.di, provideAuditRepository)
	container.MustSingleton(ki.di, provideAuditManager)
	container.MustSingleton(ki.di, provideAuditHandler)

	return nil
}
