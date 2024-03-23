package tools

import (
	"github.com/bitfield/script"
)

func Screenshot() (err error) {
	err = script.Exec("hyprshot -c -s -m output").Error()
	return
}
func NamedScreenshot(name string) (err error) {
	err = script.Exec("hyprshot -c -s -m output -f" + name).Error()
	return
}

func CheckInternetConnection() bool {
	_, err := script.Exec("ping -c 1 www.money9.com").String()
	return err == nil
}
