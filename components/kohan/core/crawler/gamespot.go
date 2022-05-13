package crawler

import (
	"github.com/PuerkitoBio/goquery"
	helper2 "github.com/amanhigh/go-fun/common/helper"
	"github.com/amanhigh/go-fun/common/util"
	crawler2 "github.com/amanhigh/go-fun/models/crawler"
	"github.com/fatih/color"
)

type GamespotCrawler struct {
	topUrl string
}

func NewGameSpotCrawler(topLink string) Crawler {
	color.Yellow("Starting Game Search")
	return &GamespotCrawler{topLink}
}

func (self *GamespotCrawler) GatherLinks(page *util.Page, ch chan crawler2.CrawlInfo) {
	games := page.Document.Find("h3.media-title")
	games.Each(func(i int, selection *goquery.Selection) {
		info := crawler2.GameInfo{Name: selection.Text()}
		if gameLink, ok := selection.Parent().Parent().Attr(util.HREF); ok {
			info.Link = helper2.GetAbsoluteLink(page, gameLink)
		}
		ch <- &info
	})
}

func (self *GamespotCrawler) NextPageLink(page *util.Page) (url string, ok bool) {
	nextPage := page.Document.Find("li.skip.next a")
	if url, ok = nextPage.Attr(util.HREF); ok {
		url = helper2.GetAbsoluteLink(page, url)
	}
	return
}

func (self *GamespotCrawler) PrintSet(good CrawlSet, bad CrawlSet) bool {
	return true
}

func (self *GamespotCrawler) GetTopPage() *util.Page {
	return util.NewPage(self.topUrl)
}
