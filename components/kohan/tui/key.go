package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HotkeyManager struct {
	uiManager *UIManager // Dependency on UIManager
}

func NewHotkeyManager(uiManager *UIManager) *HotkeyManager {
	return &HotkeyManager{
		uiManager: uiManager,
	}
}

func (h *HotkeyManager) SetupHotkeys(app *tview.Application) {
	app.SetInputCapture(h.handleHotkeys)
}

func (h *HotkeyManager) handleHotkeys(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		h.uiManager.ToggleServiceSelection()
		h.uiManager.UpdateContext()
	}
	return event
}
