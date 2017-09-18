package util

import (
	"regexp"
	"crypto/md5"
	"encoding/hex"
)

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

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}