package core

import (
	"errors"
	"fmt"
	tools2 "github.com/amanhigh/go-fun/common/tools"
	util2 "github.com/amanhigh/go-fun/common/util"
	"strings"
)

func MergeMux(splitMap map[string][]string) map[string][]string {
	muxMap := make(map[string][]string)
	for key, value := range splitMap {
		/* Collate Mux into single mux File */
		if strings.Contains(key, "mux") || strings.Contains(key, "-v1") {
			muxMap["mux"] = append(muxMap["mux"], value...)
		}
		/* Also retain all individual groups */
		muxMap[key] = value
	}
	return muxMap
}

func BuildSplitMap(lines []string) map[string][]string {
	splitMap := make(map[string][]string)
	var group string
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "["):
			group = strings.Trim(line, "[]")
			//fmt.Println("Creating New Group:", group)
			splitMap[group] = make([]string, 0)
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

func SplitAnsibleConfig(configPath string) error {
	if configPath != "" {
		lines := util2.ReadAllLines(configPath)
		splitMap := BuildSplitMap(lines)
		muxMap := MergeMux(splitMap)

		for cluster, value := range muxMap {
			if util2.IsDebugMode() {
				fmt.Println(cluster, len(value))
			}
			tools2.WriteClusterFile(cluster, strings.Join(value, "\n"))
		}
		return nil
	} else {
		return errors.New("Missing Config Path")
	}
}
