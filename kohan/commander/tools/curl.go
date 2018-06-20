package tools

import (
	"fmt"
	"strings"

	"github.com/amanhigh/go-fun/util"
)

const TIMEOUT = 10

const (
	CURL_METHOD_GET  = "GET"
	CURL_METHOD_POST = "POST"
	CURL_METHOD_PUT  = "PUT"
)

func Jcurl(url string, pipe string) (output string) {
	if util.IsDebugMode() {
		util.PrintPink(url)
	}

	if pipe == "" {
		output = Curl(url, CURL_METHOD_GET, "jq .")
	} else {
		output = Curl(url, CURL_METHOD_GET, pipe)
	}
	return
}

func Curl(url string, method string, pipe string) (output string) {
	output = RunCommandPrintError(fmt.Sprintf("curl -m %v -X%v -s '%v' | %v", TIMEOUT, method, url, pipe))
	return
}

func ContentPiperSplit(content string, pipe string) []string {
	output := ContentPiper(content, pipe)
	return util.FilterEmptyLines(strings.Split(output, "\n"))
}

func ContentPiper(content string, pipe string) string {
	output := RunCommandPrintError(fmt.Sprintf("echo '%v' | %v", content, pipe))
	return output
}
