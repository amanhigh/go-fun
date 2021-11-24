package core

import (
	"github.com/amanhigh/go-fun/apps/models/config"
	"github.com/amanhigh/go-fun/util"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
)

func DebugControl(flag bool) {
	if flag {
		color.Cyan("Enabling Debug Mode")
		ioutil.WriteFile(config.DEBUG_FILE, []byte{}, util.DEFAULT_PERM)
	} else {
		color.Red("Disabling Debug Mode")
		os.Remove(config.DEBUG_FILE)
	}
	color.Yellow("Debug Mode: %v", IsDebugMode())
}

func IsDebugMode() bool {
	return config.KOHAN_DEBUG || util.PathExists(config.DEBUG_FILE)
}
