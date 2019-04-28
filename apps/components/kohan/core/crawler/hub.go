package crawler

import (
	"github.com/PuerkitoBio/goquery"
	util2 "github.com/amanhigh/go-fun/apps/common/util"
	"github.com/amanhigh/go-fun/apps/models/crawler"
	"github.com/amanhigh/go-fun/util"
	"github.com/amanhigh/go-fun/util/helper"
)

type HubCrawler struct {
	topUrl string
}

func NewHubCrawler(topLink string) Crawler {
	util.PrintYellow("Starting Hub Search")
	return &HubCrawler{topUrl: topLink}
}

func (self *HubCrawler) GatherLinks(page *util2.Page, ch chan crawler.CrawlInfo) {
	hubs := page.Document.Find(".js-pop a")
	hubs.Each(func(i int, selection *goquery.Selection) {
		if href, ok := selection.Attr(util2.HREF); ok {
			ch <- &crawler.LinkInfo{helper.GetAbsoluteLink(page, href)}
		}
	})
}

func (self *HubCrawler) NextPageLink(page *util2.Page) (url string, ok bool) {
	nextPage := page.Document.Find(".page_next > a:nth-child(1)")
	if url, ok = nextPage.Attr(util2.HREF); ok {
		url = helper.GetAbsoluteLink(page, url)
	}
	return
}

func (self *HubCrawler) PrintSet(good []crawler.CrawlInfo, bad []crawler.CrawlInfo) bool {
	return true
}

func (self *HubCrawler) GetTopPage() *util2.Page {
	return util2.NewPage(self.topUrl)
}
