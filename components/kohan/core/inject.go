package core

import (
	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/tui"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/go-resty/resty/v2"
	"github.com/golobby/container/v3"
	"github.com/rivo/tview"
)

// Interface and implementation in same file
type KohanInterface interface {
	GetDariusApp(cfg config.DariusConfig) (*DariusV1, error)
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
	return clients.NewAlphaClient(client, ki.config.Tax.AlphaBaseURL, ki.config.Tax.AlphaAPIKey)
}

func (ki *KohanInjector) provideSBIClient(client *resty.Client) clients.SBIClient {
	return clients.NewSBIClient(client, ki.config.Tax.SBIBaseURL)
}

func (ki *KohanInjector) provideTickerManager(client clients.AlphaClient) *manager.TickerManagerImpl {
	return manager.NewTickerManager(client, ki.config.Tax.DownloadsDir)
}

func (ki *KohanInjector) provideDividendManager(sbiManager manager.SBIManager) manager.DividendManager {
	return manager.NewDividendManager(
		sbiManager,
		ki.config.Tax.DownloadsDir,
		ki.config.Tax.DividendFile,
	)
}

// Public singleton access - returns interface only
func GetKohanInterface() KohanInterface {
	return globalInjector
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
