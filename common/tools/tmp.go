package tools

import (
	"fmt"
	"io/ioutil"

	"github.com/amanhigh/go-fun/common/util"
)

const TEMP_CURL_FILE = "/tmp/curl.json"

func WriteTempCurl(data string) {
	ioutil.WriteFile(TEMP_CURL_FILE, []byte(data), util.DEFAULT_PERM)
}

func RunTempCurlCommand(cmd string) string {
	return RunCommandPrintError(getCmd(cmd))
}

func PrintTempCurlCommand(cmd string) {
	LiveCommand(getCmd(cmd))
}

func getCmd(cmd string) string {
	return fmt.Sprintf("cat %v | %v", TEMP_CURL_FILE, cmd)
}
func GetAbsoluteLink(page *util.Page, uri string) string {
	return fmt.Sprintf("https://%v%v", page.Document.Url.Host, uri)
}
