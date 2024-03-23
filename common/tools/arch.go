package tools

import (
	"strconv"
	"strings"
	"time"

	"github.com/bitfield/script"
)

func Screenshot() (err error) {
	err = script.Exec("hyprshot -c -s -m output").Error()
	return
}

func CheckInternetConnection() bool {
	_, err := script.Exec("ping -c 1 www.money9.com").String()
	return err == nil
}

func IsOSIdle(threshold time.Duration) (ok bool, err error) {
	var idleTimeMilliseconds int
	if idleTimeMilliseconds, err = IdleTimeOS(); err == nil {
		ok = int64(idleTimeMilliseconds) > threshold.Milliseconds()
	}
	return
}

func IdleTimeOS() (idleTimeMilliseconds int, err error) {
	var idleTime string
	if idleTime, err = script.Exec("xprintidle").String(); err == nil {
		idleTimeMilliseconds, err = strconv.Atoi(strings.Trim(idleTime, "\n"))
	}
	// color.Blue("Idle Time (ms): %d", idleTimeMilliseconds)
	return
}
