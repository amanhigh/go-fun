package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bitfield/script"
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
	Fullscreen int    `json:"fullscreen"`
}

type HyperMonitor struct {
	Name    string `json:"name"`
	Focused bool   `json:"focused"`
}

func GetHyperWindow() (window HyperWindow, err error) {
	var result string
	if result, err = script.Exec("hyprctl activewindow -j").String(); err == nil {
		err = json.Unmarshal([]byte(result), &window)
	}

	return
}

func GetActiveMonitor() (name string, err error) {
	var result string
	if result, err = script.Exec("hyprctl monitors -j").String(); err != nil {
		return
	}

	var monitors []HyperMonitor
	if err = json.Unmarshal([]byte(result), &monitors); err != nil {
		return
	}

	for _, m := range monitors {
		if m.Focused {
			name = m.Name
			return
		}
	}

	err = fmt.Errorf("no focused monitor found")
	return
}

func HyperDispatch(cmd string) (err error) {
	cmd = fmt.Sprintf("hyprctl dispatch %v", cmd)
	_, err = script.Exec(cmd).String()
	return
}
