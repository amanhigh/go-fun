package tools

import (
	"fmt"
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

func RunOrFocus(title string) (err error) {
	if err = FocusWindow(title); err != nil {
		_, err = RunProcess(title)
	}
	return
}

func FocusWindow(title string) (err error) {
	var windowID string
	if windowID, err = FindWindow(title); err == nil {
		err = ActivateWindow(windowID)
	} else {
		Notify("Window Not Found", title)
	}
	return
}

func FocusWindowByTitle(title string) (err error) {
	_, err = script.Exec(fmt.Sprintf("xdotool search --name '%v' windowactivate", title)).String()
	return
}

func ActivateWindow(windowId string) (err error) {
	_, err = script.Exec(fmt.Sprintf("xdotool windowactivate %v", windowId)).String()
	return
}

func FindWindow(class string) (windowId string, err error) {
	windowId, err = script.Exec(fmt.Sprintf("xdotool search --onlyvisible --class '%v'", class)).First(1).String()
	return
}

func GetActiveWindow() (windowName string, err error) {
	windowName, err = script.Exec("xdotool getactivewindow getwindowname").String()
	windowName = strings.TrimSpace(windowName)
	return
}

func GetDesktop() (number string, err error) {
	number, err = script.Exec("xdotool get_desktop").String()
	number = strings.TrimSpace(number)
	return
}

func SendKey(keys string) error {
	_, err := script.Exec(fmt.Sprintf("xdotool key --clearmodifiers %v", keys)).String()
	return err
}

func SendInput(input string) error {
	_, err := script.Exec(fmt.Sprintf("xdotool type \"%v\"", input)).String()
	return err
}
