package crawler

import (
	"context"
	"fmt"
	util2 "github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/crawler"
	"github.com/fatih/color"
	"github.com/wesovilabs/koazee"
	"io/ioutil"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	GOOD_URL_FILE = "/tmp/good.txt"
	BAD_URL_FILE  = "/tmp/bad.txt"
	BUFFER_SIZE   = 512
)

type Crawler interface {
	GatherLinks(page *util2.Page, ch chan crawler.CrawlInfo)
	NextPageLink(page *util2.Page) (string, bool)
	PrintSet(good CrawlSet, bad CrawlSet) bool
	GetTopPage() *util2.Page
}

type CrawlSet struct {
	infos []crawler.CrawlInfo
}

func (self *CrawlSet) Add(info crawler.CrawlInfo) {
	self.infos = append(self.infos, info)
}

func (self *CrawlSet) ToUrl() (urls []string) {
	for _, info := range self.infos {
		urls = append(urls, info.ToUrl()...)
	}
	urls = koazee.StreamOf(urls).RemoveDuplicates().Out().Val().([]string)
	return
}

func (self *CrawlSet) Size() (size int) {
	return len(self.ToUrl())
}

type CrawlerManager struct {
	Crawler    Crawler
	ctx        context.Context
	cancelFunc context.CancelFunc

	verbose bool

	/* Counts to track collected & required */
	collected int32
	required  int32

	infoChannel chan crawler.CrawlInfo
	goodInfo    CrawlSet
	badInfo     CrawlSet

	/* Concurrency Control */
	semaphoreChannel chan int
}

func NewCrawlerManager(crawler Crawler, requiredCount int, verbose bool) *CrawlerManager {
	return &CrawlerManager{
		Crawler:          crawler,
		required:         int32(requiredCount),
		infoChannel:      make(chan crawler.CrawlInfo, BUFFER_SIZE),
		verbose:          verbose,
		semaphoreChannel: make(chan int, runtime.NumCPU()),
	}
}

func (self *CrawlerManager) Crawl() {
	color.Yellow("Crawling RequiredLinks:%v Cores: %v", self.required, runtime.NumCPU())

	/* Fire First Crawler */
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	go self.crawlRecursive(self.Crawler.GetTopPage(), waitGroup)

	/* Collect & Organise Crawled Links */
	go self.BuildSet()

	/* Wait for all Crawlers to Return */
	waitGroup.Wait()
	close(self.infoChannel)

	/* Print Organised Links */
	self.PrintSet(self.goodInfo, self.badInfo)
}

func (self *CrawlerManager) BuildSet() {
	/* Fire Parallel Consumer to Separate Movies */
	for info := range self.infoChannel {
		if info.GoodBad() == nil {
			if util2.IsDebugMode() {
				color.Cyan("%+v", info)
			}
			self.goodInfo.Add(info)
			atomic.AddInt32(&self.collected, 1)
		} else {
			if util2.IsDebugMode() {
				color.HiMagenta("%+v", info)
			}
			self.badInfo.Add(info)
		}
	}
}

func (self *CrawlerManager) PrintSet(good CrawlSet, bad CrawlSet) {
	/* Check if Crawler want us to print or already has printed required info */
	if ok := self.Crawler.PrintSet(good, bad); ok {
		/* Output Good/Bad Info in Separate Sections */
		color.Green("Passed Info: %v", good.Size())
		self.printWriteCrawledInfo(good, GOOD_URL_FILE)

		color.Red("Failed Info: %v", bad.Size())
		self.printWriteCrawledInfo(bad, BAD_URL_FILE)
	}
}

/**
Print Info using interface and write extracted links to
GOOD/BAD Files for Chrome Processing
*/
func (self *CrawlerManager) printWriteCrawledInfo(set CrawlSet, filePath string) {
	urlDump := strings.Join(set.ToUrl(), "\n")
	if self.verbose {
		fmt.Println(urlDump)
	}
	ioutil.WriteFile(filePath, []byte(urlDump), util2.DEFAULT_PERM)
}

/**
Recursively Crawl Given Page moving to next if next link is available.
Write all Movies of current page onto channel
*/
func (self *CrawlerManager) crawlRecursive(page *util2.Page, waitGroup *sync.WaitGroup) {
	/* Aquire Grant */
	self.semaphoreChannel <- 1
	collected := atomic.LoadInt32(&self.collected)

	if collected < self.required {
		color.Yellow("Processing: %v Collected: %v", page.Document.Url.String(), collected)
		/* If Next Link is Present Crawl It */
		if link, ok := self.Crawler.NextPageLink(page); ok {
			waitGroup.Add(1)
			go self.crawlRecursive(util2.NewPage(link), waitGroup)
		}
		/* Find Links for this Page */
		self.Crawler.GatherLinks(page, self.infoChannel)
	}

	/* Release Grant */
	<-self.semaphoreChannel
	waitGroup.Done()
}
