package crawler

import (
	"github.com/PuerkitoBio/goquery"
	helper "github.com/amanhigh/go-fun/apps/common/helper"
	util2 "github.com/amanhigh/go-fun/apps/common/util"
	"github.com/amanhigh/go-fun/apps/models/crawler"
	"github.com/fatih/color"
)

type GamespotCrawler struct {
	topUrl string
}

func NewGameSpotCrawler(topLink string) Crawler {
	color.Yellow("Starting Game Search")
	return &GamespotCrawler{topLink}
}

func (self *GamespotCrawler) GatherLinks(page *util2.Page, ch chan crawler.CrawlInfo) {
	games := page.Document.Find("h3.media-title")
	games.Each(func(i int, selection *goquery.Selection) {
		info := crawler.GameInfo{Name: selection.Text()}
		if gameLink, ok := selection.Parent().Parent().Attr(util2.HREF); ok {
			info.Link = helper.GetAbsoluteLink(page, gameLink)
		}
		ch <- &info
	})
}

func (self *GamespotCrawler) NextPageLink(page *util2.Page) (url string, ok bool) {
	nextPage := page.Document.Find("li.skip.next a")
	if url, ok = nextPage.Attr(util2.HREF); ok {
		url = helper.GetAbsoluteLink(page, url)
	}
	return
}

func (self *GamespotCrawler) PrintSet(good CrawlSet, bad CrawlSet) bool {
	return true
}

func (self *GamespotCrawler) GetTopPage() *util2.Page {
	return util2.NewPage(self.topUrl)
}
