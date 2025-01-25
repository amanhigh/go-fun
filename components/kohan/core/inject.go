package core

import (
	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/tui"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/go-resty/resty/v2"
	"github.com/golobby/container/v3"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// Interface and implementation in same file
type KohanInterface interface {
	GetDariusApp(cfg config.DariusConfig) (*DariusV1, error)
	// Add new method
	GetFAManager() manager.FAManager
}

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

func (ki *KohanInjector) provideAlphaClient(client *resty.Client) clients.AlphaClient {
	return clients.NewAlphaClient(client, ki.config.FA.AlphaBaseURL, ki.config.FA.AlphaAPIKey)
}

func (ki *KohanInjector) provideSBIClient(client *resty.Client) clients.SBIClient {
	return clients.NewSBIClient(client, ki.config.FA.SBIBaseURL)
}

func (ki *KohanInjector) provideTickerManager(client clients.AlphaClient) *manager.TickerManagerImpl {
	return manager.NewTickerManager(client, ki.config.FA.DownloadsDir)
}

func (ki *KohanInjector) provideFAManager(tickerManager manager.TickerManager, sbiManager manager.SBIManager) manager.FAManager {
	return manager.NewFAManager(tickerManager, sbiManager)
}

// Public singleton access - returns interface only
func GetKohanInterface() KohanInterface {
	return globalInjector
}

func (ki *KohanInjector) GetFAManager() manager.FAManager {
	var faManager manager.FAManager
	err := ki.di.Resolve(&faManager)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get FAManager")
	}
	return faManager
}

func (ki *KohanInjector) GetDariusApp(cfg config.DariusConfig) (*DariusV1, error) {
	// Register config for this specific build
	container.MustSingleton(ki.di, func() config.DariusConfig {
		return cfg
	})

	// Register other dependencies
	container.MustSingleton(ki.di, tview.NewApplication)
	container.MustSingleton(ki.di, provideServiceManager)
	container.MustSingleton(ki.di, provideUIManager)
	container.MustSingleton(ki.di, provideHotkeyManager)

	container.MustSingleton(ki.di, ki.provideAlphaClient)
	container.MustSingleton(ki.di, ki.provideSBIClient)
	container.MustSingleton(ki.di, ki.provideTickerManager)
	container.MustSingleton(ki.di, ki.provideFAManager)

	// Build app
	app := &DariusV1{}
	err := ki.di.Fill(app)
	return app, err
}

func provideServiceManager(cfg config.DariusConfig) (serviceManager *tui.ServiceManagerImpl) {
	return tui.NewServiceManager(cfg.MakeDir, cfg.SelectedServiceFile)
}

func provideUIManager(app *tview.Application, svcManager tui.ServiceManager) *tui.UIManagerImpl {
	return tui.NewUIManager(app, svcManager)
}

func provideHotkeyManager(uiManager tui.UIManager, serviceManager tui.ServiceManager) *tui.HotkeyManagerImpl {
	return tui.NewHotkeyManager(uiManager, serviceManager)
}
