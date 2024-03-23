package tools

import (
	"fmt"

	"github.com/bitfield/script"
)

func SendKey(keys string) error {
	_, err := script.Exec(fmt.Sprintf("wtype %v", keys)).String()
	return err
}

func SendInput(input string) error {
	_, err := script.Exec(fmt.Sprintf("wtype \"%v\"", input)).String()
	return err
}

func ClipCopy(text string) (err error) {
	err = script.Echo(text).Exec("wl-copy").Error()
	return
}
