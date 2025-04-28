package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
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
	uiManager      UIManager
	serviceManager ServiceManager
	hotkeys        map[rune]Hotkey
}

func NewHotkeyManager(uiManager UIManager, serviceManager ServiceManager) *HotkeyManagerImpl {
	return &HotkeyManagerImpl{
		uiManager:      uiManager,
		serviceManager: serviceManager,
	}
}

func (h *HotkeyManagerImpl) SetupHotkeys() {
	h.uiManager.SetGlobalInputCapture(h.handleHotkeys)
	h.setupHotkeyConfig()
}

func (h *HotkeyManagerImpl) setupHotkeyConfig() {
	h.hotkeys = make(map[rune]Hotkey)
	h.setupNavigationHotkeys()
	h.setupServiceManagementHotkeys()
	h.setupSelectionHotkeys()
}

func (h *HotkeyManagerImpl) setupNavigationHotkeys() {
	hotkeys := []Hotkey{
		{Key: 'q', Description: "Quit the application", Handler: func() { h.uiManager.StopApplication() }},
		{Key: '?', Description: "Display help information", Handler: func() { h.uiManager.ShowOutput(h.GenerateHelpText()) }},
		{Key: '/', Description: "Focus on filter input", Handler: func() { h.uiManager.FocusFilterInput() }},
		{Key: 'f', Description: "Clear filter", Handler: func() {
			h.uiManager.clearFilterInput()
			h.uiManager.FocusServiceList()
		}},
	}
	h.registerHotkeys(hotkeys)
}

func (h *HotkeyManagerImpl) setupServiceManagementHotkeys() {
	hotkeys := []Hotkey{
		{Key: 'c', Description: "Clear selected services", Handler: func() { h.serviceManager.ClearSelectedServices(); h.uiManager.UpdateContext() }},
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
	h.registerHotkeys(hotkeys)
}

func (h *HotkeyManagerImpl) setupSelectionHotkeys() {
	hotkeys := []Hotkey{
		{Key: ' ', Description: "Toggle service selection or filtered services", Handler: func() {
			if h.uiManager.IsFocusOnList() {
				h.uiManager.ToggleServiceSelection()
			} else if h.uiManager.IsFocusOnFilter() {
				h.uiManager.ToggleFilteredServices()
			}
		}},
	}
	h.registerHotkeys(hotkeys)
}

func (h *HotkeyManagerImpl) registerHotkeys(hotkeys []Hotkey) {
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
	switch event.Key() { //nolint:exhaustive
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
	default:
		// Handle other keys by returning the event
		return event
	}

	return event
}
