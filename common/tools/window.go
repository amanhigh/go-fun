package tools

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func IsWindowFocused(title string) (ok bool, err error) {
	var windowTitle string
	if windowTitle, err = GetActiveWindow(); err == nil {
		// Check if the active window name contains the title case insensitive
		ok = strings.Contains(strings.ToLower(windowTitle), strings.ToLower(title))
	}
	return
}

func RunOrFocus(title string) (err error) {
	if err = FocusWindow(title); err == nil {
		var ok bool
		if ok, err = IsWindowFocused(title); err == nil && !ok {
			log.Info().Str("Window", title).Msg("Starting Process")
			_, err = RunProcess(strings.ToLower(title))
		}
	}
	return
}

func FocusWindow(title string) (err error) {
	err = HyperDispatch(fmt.Sprintf("focuswindow 'title:.*%v.*'", title))
	return
}

func FocusWorkspace(id int) (err error) {
	err = HyperDispatch(fmt.Sprintf("workspace %v", id))
	return
}

func FocusMonitor(id int) (err error) {
	err = HyperDispatch(fmt.Sprintf("focusmonitor %v", id))
	return
}

func FocusLastWindow() (err error) {
	err = HyperDispatch("focuscurrentorlast")
	return
}

func GetActiveWindow() (title string, err error) {
	var window HyperWindow
	window, err = GetHyperWindow()
	title = window.Title
	return
}
