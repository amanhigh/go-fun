package tools

import (
	"fmt"
	"strings"
)

func IsWindowFocused(title string) (ok bool, err error) {
	var windowName string
	windowName, err = GetActiveWindow()
	if err != nil {
		return
	}

	// Check if the active window name contains the title case insensitive
	ok = strings.Contains(strings.ToLower(windowName), strings.ToLower(title))

	if ok {
		Notify("Found", windowName)
	} else {
		Notify("Not Found", windowName)
	}

	return
}

func RunOrFocus(title string) (err error) {
	if err = FocusWindow(title); err != nil {
		_, err = RunProcess(title)
	}
	return
}

func FocusWindow(title string) (err error) {
	err = HyperDispatch(fmt.Sprintf("focuswindow %v", title))
	return
}

func GetActiveWindow() (title string, err error) {
	var window HyperWindow
	window, err = GetHyperWindow()
	title = window.Title
	return
}
