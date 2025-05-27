package core

import (
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/manager/tui"
	"github.com/rivo/tview"
)

type DariusV1 struct {
	app       *tview.Application `container:"type"`
	uiManager tui.UIManager      `container:"type"`
	hotkeys   tui.HotkeyManager  `container:"type"`
}

func (d *DariusV1) Run() error {
	d.uiManager.SetupLayout()
	d.hotkeys.SetupHotkeys()
	err := d.app.Run()
	return fmt.Errorf("failed to run darius app: %w", err)
}
