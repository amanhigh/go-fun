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

func (self *KohanInjector) provideAlphaClient(client *resty.Client) clients.AlphaClient {
	return clients.NewAlphaClient(client, self.config.FA.AlphaBaseURL, self.config.FA.AlphaAPIKey)
}

func (self *KohanInjector) provideSBIClient(client *resty.Client) clients.SBIClient {
	return clients.NewSBIClient(client, self.config.FA.SBIBaseURL)
}

func (self *KohanInjector) provideTickerManager(client clients.AlphaClient) *manager.TickerManagerImpl {
	return manager.NewTickerManager(client, self.config.FA.DownloadsDir)
}

func (self *KohanInjector) provideFAManager(tickerManager manager.TickerManager, sbiManager manager.SBIManager) manager.FAManager {
	return manager.NewFAManager(tickerManager, sbiManager)
}

// Public singleton access - returns interface only
func GetKohanInterface() KohanInterface {
	return globalInjector
}

func (self *KohanInjector) GetDariusApp(cfg config.DariusConfig) (*DariusV1, error) {
	// Register config for this specific build
	container.MustSingleton(self.di, func() config.DariusConfig {
		return cfg
	})

	// Register other dependencies
	container.MustSingleton(self.di, tview.NewApplication)
	container.MustSingleton(self.di, provideServiceManager)
	container.MustSingleton(self.di, provideUIManager)
	container.MustSingleton(self.di, provideHotkeyManager)

	container.MustSingleton(self.di, self.provideAlphaClient)
	container.MustSingleton(self.di, self.provideSBIClient)
	container.MustSingleton(self.di, self.provideTickerManager)
	container.MustSingleton(self.di, self.provideFAManager)

	// Build app
	app := &DariusV1{}
	err := self.di.Fill(app)
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
