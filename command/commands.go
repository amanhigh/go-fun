package commander

import "fmt"

func Jcurl(url string, pipe string) {
	PrintCommand(fmt.Sprintf("curl -s %v | jq . %v", url, pipe))
}
