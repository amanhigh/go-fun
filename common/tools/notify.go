package tools

import (
	"fmt"

	"github.com/bitfield/script"
)

func Notify(title string, message string) {
	script.Exec(fmt.Sprintf("notify-send '%v' '%v'", title, message)).Wait()
}
