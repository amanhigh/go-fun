package tools

import (
	"fmt"

	"github.com/bitfield/script"
)

func SendKey(keys string) error {
	_, err := script.Exec(fmt.Sprintf("xdotool key --clearmodifiers %v", keys)).String()
	return err
}

func SendInput(input string) error {
	_, err := script.Exec(fmt.Sprintf("xdotool type \"%v\"", input)).String()
	return err
}
