package util

import (
	"regexp"
	"strconv"
	log "github.com/Sirupsen/logrus"
)

func ReplaceRegEx(content string, search string, replace string) string {
	matcher := regexp.MustCompile(search)
	return matcher.ReplaceAllString(content, replace)
}

func FilterEmptyLines(lines []string) []string {
	nonEmptyLines := []string{}
	for _, line := range lines {
		if len(line) > 0 {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	return nonEmptyLines
}

func ParseFloat(value string) float64 {
	if float, err := strconv.ParseFloat(value, 64); err == nil {
		return float
	} else {
		log.WithFields(log.Fields{"Value": value, "Error": err}).Error("Error Parsing Float value")
		return -1
	}
}
