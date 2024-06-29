package tui

import (
	"github.com/rivo/tview"
)

/*
	TODO: Darius Improvements
	- Vim Like Keys
	- Config Files
	- Clean Selected Services
	- Minikube Control
	- Funapp Verification and Load Test.
	- New Tabs
*/

type DariusV1 struct {
	app       *tview.Application
	uiManager *UIManager
	hotkeys   *HotkeyManager
}

func NewDarius() *DariusV1 {
	app := tview.NewApplication()
	services := []string{"Mysql", "Redis", "Mongo"}
	svcManager := NewServiceManager(services)
	uiManager := NewUIManager(app, svcManager)
	return &DariusV1{
		app:       app,
		uiManager: uiManager,
		hotkeys:   NewHotkeyManager(app, uiManager),
	}
}

func (d *DariusV1) Run() error {
	d.uiManager.SetupLayout()
	d.hotkeys.SetupHotkeys()
	return d.app.Run()
}
