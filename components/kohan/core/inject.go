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

	container.MustSingleton(self.di, newServiceManager)
	container.MustSingleton(self.di, newUIManager)
	container.MustSingleton(self.di, newHotkeyManager)

	// Build App
	darius = &tui.DariusV1{}
	err = self.di.Fill(darius)
	return
}

func newServiceManager(cfg config.DariusConfig) (serviceManager *tui.ServiceManager) {
	serviceManager = &tui.ServiceManager{
		allServices:         []string{},
		selectedServices:    []string{},
		makeDir:             cfg.MakeDir,
		selectedServicePath: cfg.SelectedServiceFile,
	}
	serviceManager.loadAvailableServices()
	serviceManager.loadSelectedServices()
	return
}

func newUIManager(app *tview.Application, svcManager *tui.ServiceManager) *tui.UIManager {
	return &tui.UIManager{
		app:         app,
		mainFlex:    tview.NewFlex(),
		contextView: createTextView("Context"),
		commandView: createTextView("Command"),
		filterInput: tview.NewInputField(),
		svcManager:  svcManager,
		svcList:     createList("Services", svcManager.GetAllServices()),
	}
}

func newHotkeyManager(app *tview.Application, uiManager *tui.UIManager, serviceManager *tui.ServiceManager) *tui.HotkeyManager {
	return &tui.HotkeyManager{
		app:            app,
		uiManager:      uiManager,
		serviceManager: serviceManager,
	}
}
