package tools

import (
	"fmt"

	"github.com/bitfield/script"
	"github.com/rs/zerolog"
)

func Notify(level zerolog.Level, title, message string) (err error) {
	_, err = script.Exec(fmt.Sprintf(`hyprctl notify -1 5000 "rgb(00ff00)" "fontsize:25 %v -> %v"`, title, message)).String()
	return
}
