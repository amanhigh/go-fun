package helper

import (
	"fmt"
	"net/url"
)

func YoutubeSearch(query string) (string) {
	return fmt.Sprintf("https://www.youtube.com/results?search_query=%v", url.QueryEscape(query))
}
