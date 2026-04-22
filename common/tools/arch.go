package tools

import (
	"fmt"

	"github.com/bitfield/script"
)

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
	fullPath := dir + "/" + name
	err = script.Exec("sh -c 'grim -g \"$(slurp)\" " + fullPath + "'").Error()
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
