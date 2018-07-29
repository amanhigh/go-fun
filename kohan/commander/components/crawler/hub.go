package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/amanhigh/go-fun/models/crawler"
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

func (self *HubCrawler) GatherLinks(page *util.Page, ch chan crawler.CrawlInfo) {
	hubs := page.Document.Find(".js-pop a")
	hubs.Each(func(i int, selection *goquery.Selection) {
		if linkInfo, ok := selection.Attr(util.HREF); ok {
			ch <- &crawler.LinkInfo{linkInfo}
		}
	})
}

func (self *HubCrawler) NextPageLink(page *util.Page) (url string, ok bool) {
	nextPage := page.Document.Find(".page_next > a:nth-child(1)")
	if url, ok = nextPage.Attr(util.HREF); ok {
		url = helper.GetAbsoluteLink(page, url)
	}
	return
}

func (self *HubCrawler) PrintSet(good []crawler.CrawlInfo, bad []crawler.CrawlInfo) bool {
	return true
}

func (self *HubCrawler) GetBaseUrl() string {
	return self.topUrl
}

func (self *HubCrawler) SupplyClient() util.HttpClientInterface {
	return util.KeepAliveClient
}
