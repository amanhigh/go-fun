package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HotkeyManager struct {
	app       *tview.Application
	uiManager *UIManager
}

func (h *HotkeyManager) SetupHotkeys() {
	h.app.SetInputCapture(h.handleHotkeys)
}

func (h *HotkeyManager) handleHotkeys(event *tcell.EventKey) *tcell.EventKey {

	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'q', 'Q':
			h.app.Stop() // Quit the application
		case '?':
			h.uiManager.ShowHelp() // Display help information
		case 'c':
			h.uiManager.svcManager.ClearSelectedServices()
			h.uiManager.UpdateContext()
			return nil
		case '/':
			h.uiManager.app.SetFocus(h.uiManager.filterInput)
			return nil // Prevent slash character from being added to the input field
		case ' ':
			if h.app.GetFocus() == h.uiManager.svcList {
				h.uiManager.ToggleServiceSelection()
				h.uiManager.UpdateContext()
			} else if h.app.GetFocus() == h.uiManager.filterInput {
				h.uiManager.ToggleFilteredServices()
				return nil
			}
		}
	case tcell.KeyEsc:
		if h.app.GetFocus() == h.uiManager.filterInput {
			h.uiManager.app.SetFocus(h.uiManager.svcList)
			h.uiManager.clearFilterInput()
		}
	}

	return event
}
