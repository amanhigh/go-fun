package crawler

import (
	"fmt"
	"net/url"
)

func YoutubeSearch(query string) string {
	return fmt.Sprintf("https://www.youtube.com/results?search_query=%v", url.QueryEscape(query))
}

func YtsSearch(query string) string {
	return fmt.Sprintf("https://yts.am/browse-movies/%v/all/all/0/latest", url.QueryEscape(query))
}

func HotStarSearch(query string) string {
	return fmt.Sprintf("http://www.hotstar.com/search?q=%s", url.QueryEscape(query))
}

func TSearch(query string) string {
	return fmt.Sprintf("https://www.1337x.to/search/%s/1/", url.QueryEscape(query))
}
