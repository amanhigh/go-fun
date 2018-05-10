package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

/* PSSH */
const CLUSTER_PATH = "/tmp/clusters"
const OUTPUT_PATH = "/tmp/output"
const ERROR_PATH = "/tmp/error"
const CONSOLE_FILE = CLUSTER_PATH + "/console.txt"

const RELEASE_FILE = "/Users/amanpreet.singh/Documents/release.txt"
const DEFAULT_PARALELISM = 50

const DEBUG_FILE = "/tmp/kohandebug"

var KOHAN_DEBUG = false

func DebugControl(flag bool) {
	if flag {
		PrintSkyBlue("Enabling Debug Mode")
		ioutil.WriteFile(DEBUG_FILE, []byte{}, DEFAULT_PERM)
	} else {
		PrintRed("Disabling Debug Mode")
		os.Remove(DEBUG_FILE)
	}
	PrintYellow(fmt.Sprintf("Debug Mode: %v", IsDebugMode()))
}

func IsDebugMode() bool {
	return KOHAN_DEBUG || PathExists(DEBUG_FILE)
}
