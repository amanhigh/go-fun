package processor

import (
	"fmt"
	"flag"
	"errors"
	"github.com/amanhigh/go-fun/kohan/util"
)

type SshProcessor struct {
	Processor
}

func (p *SshProcessor) Process(commandName string) (bool) {
	var e error
	flagSet := flag.NewFlagSet(commandName, flag.ExitOnError)

	switch commandName {
	case "splitConfig":
		filePath := flagSet.String("f", "", "File Path of Ansible Config")
		flagSet.Parse(p.Args)
		e = splitAnsibleConfig(*filePath)
	default:
		fmt.Println(p.Help())
		return false
	}

	if e != nil {
		fmt.Println(e.Error())
		flagSet.Usage()
		return false
	}
	return true
}

func splitAnsibleConfig(configPath string) error {
	if configPath != "" {
		if lines, e := util.ReadLines(configPath); e == nil {
			for _, line := range lines {
				fmt.Println(line, )
			}
			return nil
		} else {
			return e
		}
	} else {
		return errors.New("Missing Config Path")
	}
}

func (p *SshProcessor) Help() string {
	return `Commands: splitConfig`
}
