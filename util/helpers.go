package util

import "regexp"

func ReplaceRegEx(content string,search string,replace string) string {
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