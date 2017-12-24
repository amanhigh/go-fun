package util

import (
	"regexp"
	"strconv"
	log "github.com/Sirupsen/logrus"
	"strings"
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

func ParseFloat(value string) (result float64) {
	var err error
	result = -1

	floatVal := strings.TrimSpace(value)
	if floatVal != "" {
		if result, err = strconv.ParseFloat(floatVal, 64); err != nil {
			log.WithFields(log.Fields{"Value": value, "Error": err}).Error("Error Parsing Float value")
		}
	}

	return
}
