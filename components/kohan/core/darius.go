package core

import (
	"github.com/amanhigh/go-fun/components/kohan/manager/tui"
	"github.com/rivo/tview"
)

type DariusV1 struct {
	app       *tview.Application     `container:"type"`
	uiManager *tui.UIManagerImpl     `container:"type"`
	hotkeys   *tui.HotkeyManagerImpl `container:"type"`
}

func (d *DariusV1) Run() error {
	d.uiManager.SetupLayout()
	d.hotkeys.SetupHotkeys()
	return d.app.Run()
}
