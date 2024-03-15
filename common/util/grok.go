package util

import (
	"regexp"
	"strings"
)

func ReplaceRegEx(content string, search string, replace string) string {
	matcher := regexp.MustCompile(search)
	return matcher.ReplaceAllString(content, replace)
}

func GoGrep(input string, pattern string) (output string) {
	compile := regexp.MustCompile(pattern)
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		submatch := compile.FindStringSubmatch(line)
		if len(submatch) > 0 {
			output += submatch[0] + "\n"
		}
	}
	return
}
