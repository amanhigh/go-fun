package processor

import (
	"fmt"
	"flag"
	"errors"
	"github.com/amanhigh/go-fun/kohan/util"
	"strings"
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
			splitMap := buildSplitMap(lines)
			muxMap := mergeMux(splitMap)

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

func mergeMux(splitMap map[string][]string) map[string][]string {
	muxMap := make(map[string][]string)
	for key, value := range splitMap {
		if strings.Contains(key, "mux") {
			muxMap["mux"] = append(muxMap["mux"], value...)
		} else {
			muxMap[key] = value
		}
	}
	return muxMap
}

func buildSplitMap(lines []string) map[string][]string {
	splitMap := make(map[string][]string)
	var group string
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "["):
			group = strings.Trim(line, "[]")
			//fmt.Println("Creating New Group:", group)
			splitMap[group] = make([]string, 1)
			break
		case strings.HasPrefix(line, "10"):
			ip := strings.Split(line, " ")[0]
			//fmt.Printf("Adding %s to %s\n", ip, group)
			splitMap[group] = append(splitMap[group], ip)
			break
		default:
		}
	}
	return splitMap
}

func (p *SshProcessor) Help() string {
	return `Commands: splitConfig`
}
