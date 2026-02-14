package core

import (
	"github.com/amanhigh/go-fun/components/kohan/manager/tui"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/golobby/container/v3"
	"github.com/rivo/tview"
)

// ---- Darius TUI Providers ----

func provideTuiServiceRepository(cfg config.DariusConfig) repository.TuiServiceRepository {
	return repository.NewTuiServiceRepository(cfg.SelectedServiceFile)
}

func provideServiceManager(cfg config.DariusConfig, repo repository.TuiServiceRepository) tui.ServiceManager {
	return tui.NewServiceManager(cfg.MakeDir, repo)
}

func provideUIManager(app *tview.Application, svcManager tui.ServiceManager) tui.UIManager {
	return tui.NewUIManager(app, svcManager)
}

func provideHotkeyManager(uiManager tui.UIManager, serviceManager tui.ServiceManager) tui.HotkeyManager {
	return tui.NewHotkeyManager(uiManager, serviceManager)
}

// ---- Darius Dependency Registration ----

// registerDariusConfig registers configuration specific to the Darius application.
func (ki *KohanInjector) registerDariusConfig(cfg config.DariusConfig) {
	// Register config for this specific build
	container.MustSingleton(ki.di, func() config.DariusConfig {
		return cfg
	})

	// Register other dependencies
	container.MustSingleton(ki.di, tview.NewApplication)
}

// registerDariusTuiManagers registers managers specific to the Darius TUI.
// Other required managers (Ticker, Account, SBI, Exchange, Tax, etc.) are registered
// by calling registerBaseDependencies, registerRepositories, registerCoreManagers, etc.
// within registerDariusDependencies.
func (ki *KohanInjector) registerDariusTuiManagers() {
	// Register TuiServiceRepository first as ServiceManager depends on it.
	container.MustSingleton(ki.di, provideTuiServiceRepository)

	// Register ServiceManager, UIManager, HotkeyManager
	container.MustSingleton(ki.di, provideServiceManager)
	container.MustSingleton(ki.di, provideUIManager)
	container.MustSingleton(ki.di, provideHotkeyManager)
}

// registerDariusDependencies registers all dependencies required for the Darius application.
func (ki *KohanInjector) registerDariusDependencies(cfg config.DariusConfig) {
	// Register base dependencies, repositories, and all managers needed for tax calculations
	// This ensures all components potentially used by Darius are available.
	ki.registerTaxDependencies() // Includes CapitalGain, Dividend, Interest, Tax managers

	// Register Darius specific configurations and TUI managers
	ki.registerDariusConfig(cfg)
	ki.registerDariusTuiManagers()
}
