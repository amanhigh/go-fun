package crawler

import (
	"github.com/amanhigh/go-fun/util"
	"context"
)

const (
	GOOD_URL_FILE = "/tmp/good.url"
	BAD_URL_FILE  = "/tmp/bad.url"
)

type Crawler interface {
	GetStartingUrl() string
	GatherLinks(page *util.Page)
	NextPageLink(page *util.Page) (string, bool)
	BuildSet()
	PrintSet()
}

type CrawlerManager struct {
	crawler    Crawler
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (self *CrawlerManager) Crawl() {
	topPage := util.NewPage(self.crawler.GetStartingUrl())

	/* Build Channel & WG for Parallel Parsers */
	//imdbInfoChannel := make(chan ImdbInfo, 512)
	//waitGroup := &sync.WaitGroup{}
	//waitGroup.Add(1)

	/* Fire First Crawler */
	self.crawlRecursive(topPage)

	/* Organise Crawled Links */
	self.crawler.BuildSet()

	/* Print Organised Links */
	self.crawler.PrintSet()
}

/**
	Recursively Crawl Given Page moving to next if next link is available.
	Write all Movies of current page onto channel
 */
func (self *CrawlerManager) crawlRecursive(page *util.Page) {
	util.PrintYellow("Processing: " + page.Document.Url.String())

	/* If Next Link is Present Crawl It */
	if link, ok := self.crawler.NextPageLink(page); ok {
		self.crawlRecursive(util.NewPage(link))
	}

	/* Find Links for this Page */
	self.crawler.GatherLinks(page)
}
