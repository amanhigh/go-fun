package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	crawler2 "github.com/amanhigh/go-fun/models/crawler"
	"github.com/fatih/color"
)

type HubCrawler struct {
	topUrl string
}

func NewHubCrawler(topLink string) Crawler {
	color.Yellow("Starting Hub Search")
	return &HubCrawler{topUrl: topLink}
}

func (self *HubCrawler) GatherLinks(page *util.Page, ch chan crawler2.CrawlInfo) {
	hubs := page.Document.Find("a[href*='video.php']")
	hubs.Each(func(i int, selection *goquery.Selection) {
		if href, ok := selection.Attr(util.HREF); ok {
			ch <- &crawler2.LinkInfo{tools.GetAbsoluteLink(page, href)}
		}
	})
}

func (self *HubCrawler) NextPageLink(page *util.Page) (url string, ok bool) {
	nextPage := page.Document.Find(".page_next > a:nth-child(1)")
	if url, ok = nextPage.Attr(util.HREF); ok {
		url = tools.GetAbsoluteLink(page, url)
	}
	return
}

func (self *HubCrawler) PrintSet(good CrawlSet, bad CrawlSet) bool {
	return true
}

func (self *HubCrawler) GetTopPage() *util.Page {
	return util.NewPage(self.topUrl)
}
