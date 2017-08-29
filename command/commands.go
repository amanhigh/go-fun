package commander

import "fmt"

func Jcurl(url string, pipe string) (output string) {
	output, _ = RunCommand(fmt.Sprintf("curl -s %v | jq . %v", url, pipe))
	return
}
