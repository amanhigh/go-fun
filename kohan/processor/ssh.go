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
			splitMap := util.BuildSplitMap(lines)
			muxMap := util.MergeMux(splitMap)

			fmt.Println("\033[1;32mAnsible Split Complete\033[0m")
			for key, value := range muxMap {
				fmt.Println(key, len(value))
				clusterPath := fmt.Sprintf("%s/%s.txt", util.CLUSTER_PATH, key)
				util.WriteLines(clusterPath, value)
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
