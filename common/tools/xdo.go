package tools

import (
	"strings"

	"github.com/bitfield/script"
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

func GetActiveWindow() (windowName string, err error) {
	windowName, err = script.Exec("xdotool getactivewindow getwindowname").String()
	return
}
