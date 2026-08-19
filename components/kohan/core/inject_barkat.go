package core

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/audit"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/golobby/container/v3"
	"github.com/rs/zerolog/log"
)

// ---- Ticker Providers ----

func provideTickerRepository(baseRepository util.BaseDbRepository) repository.TickerRepository {
	return repository.NewTickerRepository(baseRepository)
}

func provideBarkatTickerManager(repo repository.TickerRepository) manager.BarkatTickerManager {
	return manager.NewBarkatTickerManager(repo)
}

func provideTickerHandler(mgr manager.BarkatTickerManager) handler.TickerHandler {
	return handler.NewTickerHandler(mgr)
}

// ---- Alert Ticker Providers ----

func provideAlertTickerRepository(baseRepository util.BaseDbRepository) repository.AlertTickerRepository {
	return repository.NewAlertTickerRepository(baseRepository)
}

func provideAlertTickerManager(repo repository.AlertTickerRepository) manager.AlertTickerManager {
	return manager.NewAlertTickerManager(repo)
}

func provideAlertTickerHandler(mgr manager.AlertTickerManager) handler.AlertTickerHandler {
	return handler.NewAlertTickerHandler(mgr)
}

// ---- Price Alert Providers ----

func providePriceAlertRepository(baseRepository util.BaseDbRepository) repository.PriceAlertRepository {
	return repository.NewPriceAlertRepository(baseRepository)
}

func providePriceAlertManager(repo repository.PriceAlertRepository) manager.PriceAlertManager {
	return manager.NewPriceAlertManager(repo)
}

func providePriceAlertHandler(mgr manager.PriceAlertManager) handler.PriceAlertHandler {
	return handler.NewPriceAlertHandler(mgr)
}

// ---- Audit Providers ----

func provideAuditRepository(baseRepository util.BaseDbRepository) repository.AuditRepository {
	return repository.NewAuditRepository(baseRepository)
}

func provideAuditPluginRegistry(repo repository.AuditRepository) *audit.PluginRegistry {
	registry := audit.NewPluginRegistry()
	if err := registry.RegisterPlugin(audit.NewAlertCoveragePlugin(repo)); err != nil {
		log.Fatal().Err(err).Msg("failed to register audit plugin")
	}
	if err := registry.RegisterPlugin(audit.NewStaleReviewPlugin(repo)); err != nil {
		log.Fatal().Err(err).Msg("failed to register audit plugin")
	}
	return registry
}

func provideAuditManager(registry *audit.PluginRegistry) manager.AuditManager {
	return manager.NewAuditManager(registry)
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
	container.MustSingleton(ki.di, provideAuditPluginRegistry)
	container.MustSingleton(ki.di, provideAuditManager)
	container.MustSingleton(ki.di, provideAuditHandler)

	return nil
}
