package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bitfield/script"
	"github.com/rs/zerolog"
)

type HyperlandWindow struct {
	Mapped    bool `json:"mapped"`
	Workspace struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"workspace"`
	Monitor    int    `json:"monitor"`
	Class      string `json:"class"`
	Title      string `json:"title"`
	Pid        int    `json:"pid"`
	Xwayland   bool   `json:"xwayland"`
	Fullscreen bool   `json:"fullscreen"`
}

func GetActiveWindowV1() (window HyperlandWindow, err error) {
	var result string
	if result, err = script.Exec("hyprctl activewindow -j").String(); err == nil {
		err = json.Unmarshal([]byte(result), &window)
	}

	return
}

func HyperDispatch(cmd string) (err error) {
	cmd = fmt.Sprintf("hyprctl dispatch %v", cmd)
	_, err = script.Exec(cmd).String()
	return
}

func NotifyV1(level zerolog.Level, message string) (err error) {
	_, err = script.Exec(fmt.Sprintf(`hyprctl notify -1 5000 "rgb(00ff00)" "fontsize:25 %v"`, message)).String()
	return
}
