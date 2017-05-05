package processor

import (
	"fmt"
	"flag"
	"errors"
)

type SshProcessor struct {
	Processor
}

func (p *SshProcessor) Process(commandName string) (bool) {
	var e error
	switch commandName {
	case "splitConfig":
		filePath := flag.String("f", "", "File Path of Ansible Config")
		flag.Parse()
		e = splitAnsibleConfig(*filePath)
	case "help":
		fmt.Println(p.Help())
	default:
		e = errors.New("Unknown Command: " + commandName)
		fmt.Println(p.Help())
	}

	if e != nil {
		print(e)
		return false
	}
	return true
}

func splitAnsibleConfig(configPath string) error {
	return nil
}

func (p *SshProcessor) Help() string {
	return `Commands: splitConfig`
}
