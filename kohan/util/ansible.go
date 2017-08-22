package util

import "strings"

func MergeMux(splitMap map[string][]string) map[string][]string {
	muxMap := make(map[string][]string)
	for key, value := range splitMap {
		/* Collate Mux into single mux File */
		if strings.Contains(key, "mux") || strings.Contains(key,"-g1-") {
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