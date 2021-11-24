package tools

import (
	"fmt"
	util "github.com/amanhigh/go-fun/apps/common/util"
	"github.com/amanhigh/go-fun/apps/components/kohan/core"
	"github.com/fatih/color"
	"strings"
)

const TIMEOUT = 60

const (
	CURL_METHOD_GET  = "GET"
	CURL_METHOD_POST = "POST"
	CURL_METHOD_PUT  = "PUT"
)

func Jcurl(url string, pipe string) (output string) {
	if core.IsDebugMode() {
		color.Magenta(url)
	}

	if pipe == "" {
		output = CurlGet(url, "jq .")
	} else {
		output = CurlGet(url, pipe)
	}
	return
}

func CurlGet(url string, pipe string) (output string) {
	output = Curl(url, CURL_METHOD_GET, "", pipe)
	return
}

func CurlPut(url string, filePath string, params string, pipe string) (output string) {
	output = Curl(url, CURL_METHOD_PUT, fmt.Sprintf("-d @%v %v", filePath, params), pipe)
	return
}

func Curl(url string, method string, params string, pipe string) (output string) {
	cmd := fmt.Sprintf("curl -m %v -X%v -s '%v' %v", TIMEOUT, method, url, params)
	if pipe != "" {
		cmd += " | " + pipe
	}
	output = RunCommandPrintError(cmd)
	return
}

func ContentPiperSplit(content string, pipe string) []string {
	output := ContentPiper(content, pipe)
	lines := strings.Split(output, "\n")
	return util.FilterEmptyLines(lines)
}

func ContentPiper(content string, pipe string) string {
	output := RunCommandPrintError(fmt.Sprintf("echo '%v' | %v", content, pipe))
	return output
}
