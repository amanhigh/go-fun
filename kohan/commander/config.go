package commander

import (
	"os"
	"io/ioutil"
	"fmt"
	"github.com/amanhigh/go-fun/util"
)

/* PSSH */
const CLUSTER_PATH = "/tmp/clusters"
const OUTPUT_PATH = "/tmp/output"
const ERROR_PATH = "/tmp/error"
const CONSOLE_FILE = CLUSTER_PATH + "/console.txt"

const RELEASE_FILE = "/Users/amanpreet.singh/Documents/release.txt"
const DEFAULT_PARALELISM = 50

const DEBUG_FILE string = "/tmp/kohandebug"

/* Date/Time */
const LAYOUT = "Jan 2, 2006 at 3:04pm (MST)"

func DebugControl(flag bool) {
	if flag {
		util.PrintSkyBlue("Enabling Debug Mode")
		ioutil.WriteFile(DEBUG_FILE, []byte{}, util.DEFAULT_PERM)
	} else {
		util.PrintRed("Disabling Debug Mode")
		os.Remove(DEBUG_FILE)
	}
	util.PrintYellow(fmt.Sprintf("Debug Mode: %v", IsDebugMode()))
}

func IsDebugMode() bool {
	return util.PathExists(DEBUG_FILE)
}
