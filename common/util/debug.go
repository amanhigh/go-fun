package util

import (
	"os"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/fatih/color"
)

func DebugControl(flag bool) {
	if flag {
		color.Cyan("Enabling Debug Mode")
		os.WriteFile(config.DEBUG_FILE, []byte{}, DEFAULT_PERM)
	} else {
		color.Red("Disabling Debug Mode")
		os.Remove(config.DEBUG_FILE)
	}
	color.Yellow("Debug Mode: %v", IsDebugMode())
}

func IsDebugMode() bool {
	return config.KOHAN_DEBUG || PathExists(config.DEBUG_FILE)
}
