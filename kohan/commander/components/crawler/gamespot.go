package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/amanhigh/go-fun/models/crawler"
	"github.com/amanhigh/go-fun/util"
	"github.com/amanhigh/go-fun/util/helper"
)

type GamespotCrawler struct {
	topUrl string
}

func NewGameSpotCrawler(topLink string) Crawler {
	util.PrintYellow("Starting Game Search")
	return &GamespotCrawler{topLink}
}

func (self *GamespotCrawler) GatherLinks(page *util.Page, ch chan crawler.CrawlInfo) {
	games := page.Document.Find("h3.media-title")
	games.Each(func(i int, selection *goquery.Selection) {
		info := crawler.GameInfo{Name: selection.Text()}
		if gameLink, ok := selection.Parent().Parent().Attr(util.HREF); ok {
			info.Link = helper.GetAbsoluteLink(page, gameLink)
		}
		ch <- &info
	})
}

func (self *GamespotCrawler) NextPageLink(page *util.Page) (url string, ok bool) {
	nextPage := page.Document.Find("li.skip.next a")
	if url, ok = nextPage.Attr(util.HREF); ok {
		url = helper.GetAbsoluteLink(page, url)
	}
	return
}

func (self *GamespotCrawler) PrintSet(good []crawler.CrawlInfo, bad []crawler.CrawlInfo) bool {
	return true
}

func (self *GamespotCrawler) GetTopPage() *util.Page {
	return util.NewPage(self.topUrl)
}
