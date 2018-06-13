package crawler

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/amanhigh/go-fun/models/crawler"
	"github.com/amanhigh/go-fun/util"
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
			info.Link = getGamespotLink(page, gameLink)
		}
		ch <- &info
	})
}

func (self *GamespotCrawler) NextPageLink(page *util.Page) (url string, ok bool) {
	nextPage := page.Document.Find("li.skip.next a")
	if url, ok = nextPage.Attr(util.HREF); ok {
		url = getGamespotLink(page, url)
	}
	return
}

func (self *GamespotCrawler) PrintSet(good []crawler.CrawlInfo, bad []crawler.CrawlInfo) bool {
	return true
}

func (self *GamespotCrawler) GetBaseUrl() string {
	return self.topUrl
}
func (self *GamespotCrawler) SupplyClient() util.HttpClientInterface {
	return util.KeepAliveClient
}

func getGamespotLink(page *util.Page, uri string) string {
	return fmt.Sprintf("https://%v%v", page.Document.Url.Host, uri)
}
