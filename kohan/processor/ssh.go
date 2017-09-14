package processor

import (
	"fmt"
	"flag"
	"errors"
	"github.com/amanhigh/go-fun/kohan/util"
)

type SshProcessor struct {
}

func (self *SshProcessor) GetHandleMap() (map[string]HandleFunc) {
	return map[string]HandleFunc{
		"splitConfig": self.handleSplitConfig,
	}
}

func (self *SshProcessor) handleSplitConfig(flagSet *flag.FlagSet, args []string) error {
	filePath := flagSet.String("f", "", "File Path of Ansible Config")
	flagSet.Parse(args)
	return SplitAnsibleConfig(*filePath)
}

func SplitAnsibleConfig(configPath string) error {
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