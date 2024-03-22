package tools

import (
	"fmt"

	"github.com/bitfield/script"
	"github.com/rs/zerolog"
)

func Notify(title string, message string) {
	script.Exec(fmt.Sprintf("notify-send '%v' '%v'", title, message)).Wait()
}

func NotifyV1(level zerolog.Level, message string) (err error) {
	_, err = script.Exec(fmt.Sprintf(`hyprctl notify -1 5000 "rgb(00ff00)" "fontsize:25 %v"`, message)).String()
	return
}
