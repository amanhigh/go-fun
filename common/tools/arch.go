package tools

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bitfield/script"
)

// ErrScreenshotAborted is returned when the user cancels a region screenshot
// (e.g., pressing Escape during slurp region selection).
var ErrScreenshotAborted = errors.New("screenshot aborted")

func Screenshot() (err error) {
	var monitor string
	if monitor, err = GetActiveMonitor(); err != nil {
		return
	}
	err = script.Exec(fmt.Sprintf("grim -o %s - | wl-copy", monitor)).Error()
	return
}

func NamedScreenshot(dir, name string) (err error) {
	var monitor string
	if monitor, err = GetActiveMonitor(); err != nil {
		return
	}
	fullPath := dir + "/" + name
	err = script.Exec(fmt.Sprintf("grim -o %s %s", monitor, fullPath)).Error()
	return
}

func NamedRegionScreenshot(dir, name string) (err error) {
	// Step 1: Run slurp to get the selected region geometry.
	// slurp exits with code 1 when the user cancels (Escape).
	geometry, err := script.Exec("slurp").String()
	if err != nil {
		return ErrScreenshotAborted
	}

	// Step 2: Use the geometry string for grim capture.
	fullPath := dir + "/" + name
	geometry = strings.TrimSpace(geometry)
	err = script.Exec(fmt.Sprintf("grim -g %s %s", geometry, fullPath)).Error()
	return
}

func CheckInternetConnection() bool {
	_, err := script.Exec("ping -c 1 www.money9.com").String()
	return err == nil
}

func PromptText(text string) (result string, err error) {
	result, err = script.Echo(text).Exec("zenity --editable --text-info").String()
	return
}
