package commander

import "fmt"

func Jcurl(url string, pipe string) (output string) {
	output = RunCommandPrintError(fmt.Sprintf("curl -s %v | jq . %v", url, pipe))
	return
}

func ContentPiper(content string,pipe string) (output string)  {
	output = RunCommandPrintError(fmt.Sprintf("echo '%v' | %v",content,pipe ))
	return
}