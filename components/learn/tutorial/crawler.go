package tutorial

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
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
	Crawl(site, 4, fetcher, urlMap)
	return
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, urlMap SafeMap) {
	//fmt.Println("Crawl Recived:", url, depth)
	// XXX: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Debug("Fetch Fail")
		return
	}

	log.WithFields(log.Fields{"Url": url, "Title": body}).Debug("URL_HIT")

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(urls))
	for _, url := range urls {
		//fmt.Println("Submitting To Crawler:", url, depth - 1)
		go func(u string) {
			defer waitGroup.Done()
			if !urlMap.Contains(u) {
				urlMap.Add(u)
				Crawl(u, depth-1, fetcher, urlMap)
			}
		}(url)
	}
	waitGroup.Wait()

	return
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
