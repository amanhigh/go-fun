package tutorial

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type SafeMap struct {
	m   map[string]bool
	mux sync.Mutex
}

func (m *SafeMap) Add(url string) {
	m.mux.Lock()
	m.m[url] = true
	m.mux.Unlock()
}

func (m *SafeMap) Contains(url string) (ok bool) {
	m.mux.Lock()
	_, ok = m.m[url]
	m.mux.Unlock()
	return
}

func StartCrawl(site string) (urlMap SafeMap) {
	/** Seed UrlMap With Top Url */
	urlMap = SafeMap{m: map[string]bool{site: true}}
	Crawl(site, 4, fetcher, &urlMap)
	return
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, urlMap *SafeMap) {
	log.Debug().Str("Url", url).Int("Depth", depth).Msg("CRAWL_RECIVED")
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		log.Debug().Err(err).Msg("Fetch Fail")
		return
	}

	log.Debug().Str("Url", url).Str("Title", body).Msg("URL_HIT")

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(urls))
	for _, url := range urls {
		log.Debug().Str("Url", url).Int("Depth", depth-1).Msg("CRAWL_SUBMIT")
		go func(u string) {
			defer waitGroup.Done()
			if !urlMap.Contains(u) {
				urlMap.Add(u)
				Crawl(u, depth-1, fetcher, urlMap)
			}
		}(url)
	}
	waitGroup.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	/** Nice Syntax for if key exists do something on value */
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}

	/** This is how you Format Error (Not Printf) */
	return "", nil, fmt.Errorf("URL_MISS: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
