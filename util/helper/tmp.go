package helper

import (
	"fmt"
	"io/ioutil"

	"github.com/amanhigh/go-fun/kohan/commander/tools"
	"github.com/amanhigh/go-fun/util"
)

const TEMP_CURL_FILE = "/tmp/curl.json"

func WriteTempCurl(data string) {
	ioutil.WriteFile(TEMP_CURL_FILE, []byte(data), util.DEFAULT_PERM)
}

func RunTempCurlCommand(cmd string) string {
	return tools.RunCommandPrintError(getCmd(cmd))
}

func PrintTempCurlCommand(cmd string) {
	tools.LiveCommand(getCmd(cmd))
}

func getCmd(cmd string) string {
	return fmt.Sprintf("cat %v | %v", TEMP_CURL_FILE, cmd)
}
