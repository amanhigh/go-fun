package processor

import (
	"fmt"
)

type SshProcessor struct {
	Processor
}

func (p *SshProcessor) Process(commandName string) (bool) {
	switch commandName {
	default:
		fmt.Println("Unknown Command:", commandName)
		return false
	}
	return true
}

func splitAnsibleConfig(configPath string) {

}
