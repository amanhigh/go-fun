package commander

import (
	"fmt"
	"strings"
)

func Jcurl(url string, pipe string) (output string) {
	if pipe == "" {
		output = RunCommandPrintError(fmt.Sprintf("curl -s '%v' | jq .", url))
	} else {
		output = RunCommandPrintError(fmt.Sprintf("curl -s '%v' | jq . | %v", url, pipe))
	}
	return
}

func ContentPiperSplit(content string, pipe string) ([]string) {
	output := ContentPiper(content, pipe)
	return FilterEmptyLines(strings.Split(output, "\n"))
}

func ContentPiper(content string, pipe string) (string) {
	output := RunCommandPrintError(fmt.Sprintf("echo '%v' | %v", content, pipe))
	return output
}
