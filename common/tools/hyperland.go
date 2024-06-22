package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitfield/script"
	"github.com/rs/zerolog/log"
)

type HyperWindow struct {
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

var (
	isSubMapActive = false
)

func GetHyperWindow() (window HyperWindow, err error) {
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

// ActivateSubmap is a Go function to activate a submap based on the window title.
//
// It takes two parameters: submap (string) and windowTitle (string) and returns an error.
func ActivateSubmap(submap, windowTitle string) (err error) {
	var window HyperWindow
	window, err = GetHyperWindow()
	if err != nil {
		return
	}
	windowMatch := strings.Contains(window.Title, windowTitle)

	if !isSubMapActive && windowMatch {
		log.Info().Str("Window", window.Title).Err(err).Msg("Enable Submap")
		isSubMapActive = true
		err = HyperDispatch("submap " + submap)
	}

	if isSubMapActive && !windowMatch {
		log.Info().Str("Window", window.Title).Err(err).Msg("Disable Submap")
		err = HyperDispatch("submap reset")
		isSubMapActive = false
	}

	return
}
