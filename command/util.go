package commander

import "regexp"

func ReplaceRegEx(content string,search string,replace string) string {
	matcher := regexp.MustCompile(search)
	return matcher.ReplaceAllString(content, replace)
}