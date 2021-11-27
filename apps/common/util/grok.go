package util

import (
	"github.com/thoas/go-funk"
	"regexp"
	"strings"
)

func ReplaceRegEx(content string, search string, replace string) string {
	matcher := regexp.MustCompile(search)
	return matcher.ReplaceAllString(content, replace)
}

func FilterEmptyLines(lines []string) []string {
	return funk.FilterString(lines, func(s string) bool {
		return len(s) > 0
	})
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
