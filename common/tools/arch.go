package tools

import (
	"github.com/bitfield/script"
)

func Screenshot() (err error) {
	err = script.Exec("hyprshot -c -s -m output").Error()
	return
}
func NamedScreenshot(dir, name string) (err error) {
	err = script.Exec("hyprshot -c -s -m output -o " + dir + " -f" + name).Error()
	return
}

func NamedRegionScreenshot(name string) (err error) {
	err = script.Exec("hyprshot -s -m region -f" + name).Error()
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
