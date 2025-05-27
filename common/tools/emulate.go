package tools

import (
	"fmt"

	"github.com/bitfield/script"
)

func SendKey(keys string) error {
	_, err := script.Exec(fmt.Sprintf("wtype %v", keys)).String()
	if err != nil {
		return fmt.Errorf("failed to send key: %w", err)
	}
	return nil
}

func SendInput(input string) error {
	_, err := script.Exec(fmt.Sprintf("wtype \"%v\"", input)).String()
	if err != nil {
		return fmt.Errorf("failed to send input: %w", err)
	}
	return nil
}

func ClipCopy(text string) (err error) {
	err = script.Echo(text).Exec("wl-copy").Error()
	return
}

func ClipPaste() (text string, err error) {
	text, err = script.Exec("wl-paste -n").String()
	return
}
