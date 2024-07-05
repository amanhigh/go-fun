package tui

import (
	"github.com/rivo/tview"
)

/*
	TODO: Darius Improvements
	- Vim Like Keys
	- Config Files
	- Minikube Control
	- Funapp Verification and Load Test.
	- New Tabs
*/

type DariusV1 struct {
	app       *tview.Application `container:"type"`
	uiManager *UIManager         `container:"type"`
	hotkeys   *HotkeyManager     `container:"type"`
}

func (d *DariusV1) Run() error {
	d.uiManager.SetupLayout()
	d.hotkeys.SetupHotkeys()
	return d.app.Run()
}
