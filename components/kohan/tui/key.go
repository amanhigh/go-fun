package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Hotkey struct {
	Key         rune
	Description string
	Handler     func()
}

type HotkeyManager struct {
	app            *tview.Application
	uiManager      *UIManager
	serviceManager *ServiceManager
	hotkeys        map[rune]Hotkey
}

func (h *HotkeyManager) SetupHotkeys() {
	h.app.SetInputCapture(h.handleHotkeys)
	h.setupHotkeyConfig()
}

func (h *HotkeyManager) setupHotkeyConfig() {
	h.hotkeys = make(map[rune]Hotkey)
	hotkeys := []Hotkey{
		{Key: 'q', Description: "Quit the application", Handler: func() { h.app.Stop() }},
		{Key: '?', Description: "Display help information", Handler: func() { h.uiManager.ShowOutput(h.GenerateHelpText()) }},
		{Key: 'c', Description: "Clear selected services", Handler: func() { h.serviceManager.ClearSelectedServices(); h.uiManager.UpdateContext() }},
		{Key: '/', Description: "Focus on filter input", Handler: func() { h.uiManager.app.SetFocus(h.uiManager.filterInput) }},
		{Key: ' ', Description: "Toggle service selection or filtered services", Handler: func() {
			if h.app.GetFocus() == h.uiManager.svcList {
				h.uiManager.ToggleServiceSelection()
				h.uiManager.UpdateContext()
			} else if h.app.GetFocus() == h.uiManager.filterInput {
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
	}
	for _, hotkey := range hotkeys {
		h.hotkeys[hotkey.Key] = hotkey
	}
}

func (h *HotkeyManager) GenerateHelpText() string {
	helpText := "Help:\n"
	for _, hotkey := range h.hotkeys {
		helpText += fmt.Sprintf("- %c: %s\n", hotkey.Key, hotkey.Description)
	}
	helpText += "\n"
	return helpText
}
func (h *HotkeyManager) handleHotkeys(event *tcell.EventKey) *tcell.EventKey {

	switch event.Key() {
	case tcell.KeyRune:
		if hotkey, exists := h.hotkeys[event.Rune()]; exists && h.app.GetFocus() != h.uiManager.filterInput {
			hotkey.Handler()
			return nil
		}
	case tcell.KeyEsc:
		if h.app.GetFocus() == h.uiManager.filterInput {
			h.uiManager.app.SetFocus(h.uiManager.svcList)
			h.uiManager.clearFilterInput()
		}
	}

	return event
}
