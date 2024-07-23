package tui

import (
	"github.com/rivo/tview"
)

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
