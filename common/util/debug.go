package util

import (
	config2 "github.com/amanhigh/go-fun/models/config"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
)

func DebugControl(flag bool) {
	if flag {
		color.Cyan("Enabling Debug Mode")
		ioutil.WriteFile(config2.DEBUG_FILE, []byte{}, DEFAULT_PERM)
	} else {
		color.Red("Disabling Debug Mode")
		os.Remove(config2.DEBUG_FILE)
	}
	color.Yellow("Debug Mode: %v", IsDebugMode())
}

func IsDebugMode() bool {
	return config2.KOHAN_DEBUG || PathExists(config2.DEBUG_FILE)
}
