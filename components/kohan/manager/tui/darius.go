package tui

import (
	"github.com/rivo/tview"
)

type DariusV1 struct {
	app       *tview.Application `container:"type"`
	uiManager *UIManagerImpl     `container:"type"`
	hotkeys   *HotkeyManagerImpl `container:"type"`
}

// BUG: #C Move to Core
func (d *DariusV1) Run() error {
	d.uiManager.SetupLayout()
	d.hotkeys.SetupHotkeys()
	return d.app.Run()
}
