package core

import (
	"github.com/amanhigh/go-fun/components/kohan/manager/tui"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/golobby/container/v3"
	"github.com/rivo/tview"
)

// FIXME: Upgrade to Kohan Injector ?
type DariusInjector struct {
	di     container.Container
	config config.DariusConfig
}

func NewDariusInjector(cfg config.DariusConfig) (di *DariusInjector) {
	return &DariusInjector{container.New(), cfg}
}

func (self *DariusInjector) BuildApp() (darius *tui.DariusV1, err error) {
	container.MustSingleton(self.di, func() config.DariusConfig {
		return self.config
	})

	container.MustSingleton(self.di, tview.NewApplication)

	container.MustSingleton(self.di, provideServiceManager)
	container.MustSingleton(self.di, provideUIManager)
	container.MustSingleton(self.di, provideHotkeyManager)

	// Build App
	darius = &tui.DariusV1{}
	err = self.di.Fill(darius)
	return
}

func provideServiceManager(cfg config.DariusConfig) (serviceManager *tui.ServiceManager) {
	return tui.NewServiceManager(cfg.MakeDir, cfg.SelectedServiceFile)
}

func provideUIManager(app *tview.Application, svcManager tui.ServiceManagerInterface) *tui.UIManager {
	return tui.NewUIManager(app, svcManager)
}

func provideHotkeyManager(app *tview.Application, uiManager tui.UIManagerInterface, serviceManager tui.ServiceManagerInterface) *tui.HotkeyManager {
	return tui.NewHotkeyManager(app, uiManager, serviceManager)
}
