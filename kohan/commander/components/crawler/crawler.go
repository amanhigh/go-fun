package crawler

import (
	"github.com/amanhigh/go-fun/util"
	"sync"
)

const (
	GOOD_URL_FILE = "/tmp/good.url"
	BAD_URL_FILE  = "/tmp/bad.url"
)

type Crawler interface {
	GetBaseUrl() string
	GatherLinks(page *util.Page)
	NextPageLink(page *util.Page) (string, bool)
	GatherComplete()
	BuildSet()
	PrintSet()
}

type CrawlerManager struct {
	Crawler Crawler
}

func (self *CrawlerManager) Crawl() {
	topPage := util.NewPage(self.Crawler.GetBaseUrl())

	/* Fire First Crawler */
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	go self.crawlRecursive(topPage, waitGroup)

	/* Collect & Organise Crawled Links */
	go self.Crawler.BuildSet()

	/* Wait for all Crawlers to Return */
	waitGroup.Wait()
	self.Crawler.GatherComplete()

	/* Print Organised Links */
	self.Crawler.PrintSet()
}

/**
	Recursively Crawl Given Page moving to next if next link is available.
	Write all Movies of current page onto channel
 */
func (self *CrawlerManager) crawlRecursive(page *util.Page, waitGroup *sync.WaitGroup) {
	util.PrintYellow("Processing: " + page.Document.Url.String())

	/* If Next Link is Present Crawl It */
	if link, ok := self.Crawler.NextPageLink(page); ok {
		waitGroup.Add(1)
		go self.crawlRecursive(util.NewPage(link), waitGroup)
	}

	/* Find Links for this Page */
	self.Crawler.GatherLinks(page)
	waitGroup.Done()
}
