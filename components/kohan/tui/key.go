package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HotkeyManager struct {
	app       *tview.Application
	uiManager *UIManager
}

func NewHotkeyManager(app *tview.Application, uiManager *UIManager) *HotkeyManager {
	return &HotkeyManager{
		app:       app,
		uiManager: uiManager,
	}
}

func (h *HotkeyManager) SetupHotkeys() {
	h.app.SetInputCapture(h.handleHotkeys)
}

func (h *HotkeyManager) handleHotkeys(event *tcell.EventKey) *tcell.EventKey {

	switch event.Key() {
	case tcell.KeyEnter:
		h.uiManager.ToggleServiceSelection()
		h.uiManager.UpdateContext()
	}
	return event
}
