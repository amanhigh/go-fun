package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HotkeyManager interface {
	SetupHotkeys()
	GenerateHelpText() string
}

type Hotkey struct {
	Key         rune
	Description string
	Handler     func()
}

type HotkeyManagerImpl struct {
	// XXX: Completely Remove App Depenency ?
	app            *tview.Application
	uiManager      UIManager
	serviceManager ServiceManager
	hotkeys        map[rune]Hotkey
}

func NewHotkeyManager(app *tview.Application, uiManager UIManager, serviceManager ServiceManager) *HotkeyManagerImpl {
	return &HotkeyManagerImpl{
		app:            app,
		uiManager:      uiManager,
		serviceManager: serviceManager,
	}
}

func (h *HotkeyManagerImpl) SetupHotkeys() {
	h.app.SetInputCapture(h.handleHotkeys)
	h.setupHotkeyConfig()
}

func (h *HotkeyManagerImpl) setupHotkeyConfig() {
	h.hotkeys = make(map[rune]Hotkey)
	hotkeys := []Hotkey{
		{Key: 'q', Description: "Quit the application", Handler: func() { h.app.Stop() }},
		{Key: '?', Description: "Display help information", Handler: func() { h.uiManager.ShowOutput(h.GenerateHelpText()) }},
		{Key: 'c', Description: "Clear selected services", Handler: func() { h.serviceManager.ClearSelectedServices(); h.uiManager.UpdateContext() }},
		{Key: '/', Description: "Focus on filter input", Handler: func() { h.uiManager.FocusFilterInput() }},
		{Key: ' ', Description: "Toggle service selection or filtered services", Handler: func() {
			if h.uiManager.IsFocusOnList() {
				h.uiManager.ToggleServiceSelection()
			} else if h.uiManager.IsFocusOnFilter() {
				h.uiManager.ToggleFilteredServices()
			}
		}},
		{Key: 's', Description: "Setup services", Handler: func() {
			output, err := h.serviceManager.SetupServices()
			if err != nil {
				h.uiManager.ShowError(err)
			} else {
				h.uiManager.ShowOutput(output)
			}
		}},
		{Key: 'C', Description: "Clean services", Handler: func() {
			output, err := h.serviceManager.CleanServices()
			if err != nil {
				h.uiManager.ShowError(err)
			} else {
				h.uiManager.ShowOutput(output)
			}
		}},
		{Key: 'u', Description: "Update services", Handler: func() {
			output, err := h.serviceManager.UpdateServices()
			if err != nil {
				h.uiManager.ShowError(err)
			} else {
				h.uiManager.ShowOutput(output)
			}
		}},
		{Key: 'f', Description: "Clear filter", Handler: func() {
			h.uiManager.clearFilterInput()
			h.uiManager.FocusServiceList()
		}},
	}
	for _, hotkey := range hotkeys {
		h.hotkeys[hotkey.Key] = hotkey
	}
}

func (h *HotkeyManagerImpl) GenerateHelpText() string {
	helpText := "Help:\n"
	for _, hotkey := range h.hotkeys {
		helpText += fmt.Sprintf("- %c: %s\n", hotkey.Key, hotkey.Description)
	}
	helpText += "\n"
	return helpText
}

func (h *HotkeyManagerImpl) handleHotkeys(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		if hotkey, exists := h.hotkeys[event.Rune()]; exists {
			// Handle space key specially
			if event.Rune() == ' ' {
				hotkey.Handler()
				return nil
			}
			// For other hotkeys, only handle if not in filter input
			if !h.uiManager.IsFocusOnFilter() {
				hotkey.Handler()
				return nil
			}
		}
	case tcell.KeyEsc:
		if h.uiManager.IsFocusOnFilter() {
			h.uiManager.FocusServiceList()
			return nil
		}
	}

	return event
}
