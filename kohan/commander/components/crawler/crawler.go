package crawler

import (
	"github.com/amanhigh/go-fun/util"
	"sync"
	"context"
	"sync/atomic"
	"fmt"
)

const (
	GOOD_URL_FILE = "/tmp/good.url"
	BAD_URL_FILE  = "/tmp/bad.url"
)

type Crawler interface {
	GetBaseUrl() string
	GatherLinks(page *util.Page) int
	NextPageLink(page *util.Page) (string, bool)
	GatherComplete()
	BuildSet()
	PrintSet()
}

type CrawlerManager struct {
	Crawler    Crawler
	ctx        context.Context
	cancelFunc context.CancelFunc

	collectCount  int32
	RequiredCount int32
}

func NewCrawlerManager(crawler Crawler, requiredCount int) *CrawlerManager {
	return &CrawlerManager{
		Crawler:       crawler,
		RequiredCount: int32(requiredCount),
	}
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
	collected := atomic.LoadInt32(&self.collectCount)
	if collected < self.RequiredCount {
		util.PrintYellow(fmt.Sprintf("Processing: %v Collected: %v", page.Document.Url.String(), collected))

		/* If Next Link is Present Crawl It */
		if link, ok := self.Crawler.NextPageLink(page); ok {
			waitGroup.Add(1)
			go self.crawlRecursive(util.NewPage(link), waitGroup)
		}
		/* Find Links for this Page */
		linksGathered := self.Crawler.GatherLinks(page)
		atomic.AddInt32(&self.collectCount, int32(linksGathered))
	}
	waitGroup.Done()
}
