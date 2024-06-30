package tui

import (
	"github.com/amanhigh/go-fun/models/config"
	"github.com/golobby/container/v3"
	"github.com/rivo/tview"
)

type DariusInjector struct {
	di     container.Container
	config config.DariusConfig
}

func NewDariusInjector(cfg config.DariusConfig) (di *DariusInjector) {
	return &DariusInjector{container.New(), cfg}
}

func (self *DariusInjector) BuildApp() (darius *DariusV1, err error) {
	container.MustSingleton(self.di, func() config.DariusConfig {
		return self.config
	})

	container.MustSingleton(self.di, tview.NewApplication)

	container.MustSingleton(self.di, newServiceManager)
	container.MustSingleton(self.di, newUIManager)
	container.MustSingleton(self.di, newHotkeyManager)

	// Build App
	darius = &DariusV1{}
	err = self.di.Fill(darius)
	return
}

func newServiceManager(cfg config.DariusConfig) (serviceManager *ServiceManager) {
	serviceManager = &ServiceManager{
		allServices:         []string{},
		selectedServices:    []string{},
		makeDir:             cfg.MakeDir,
		selectedServicePath: cfg.SelectedServiceFile,
	}
	serviceManager.loadSelectedServices()
	serviceManager.loadAvailableServices()
	return
}

func newUIManager(app *tview.Application, svcManager *ServiceManager) *UIManager {
	return &UIManager{
		app:         app,
		mainFlex:    tview.NewFlex(),
		contextView: createTextView("Context"),
		commandView: createTextView("Command"),
		filterInput: tview.NewInputField(),
		svcManager:  svcManager,
		svcList:     createList("Services", svcManager.GetAllServices()),
	}
}

func newHotkeyManager(app *tview.Application, uiManager *UIManager, serviceManager *ServiceManager) *HotkeyManager {
	return &HotkeyManager{
		app:            app,
		uiManager:      uiManager,
		serviceManager: serviceManager,
	}
}
