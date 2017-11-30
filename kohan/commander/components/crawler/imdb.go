package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"github.com/amanhigh/go-fun/util"
	. "github.com/amanhigh/go-fun/models/crawler"
	log "github.com/Sirupsen/logrus"
)

type ImdbCrawler struct {
	cutoff int
	topUrl string
}

func NewImdbCrawler(year int, language string, cutoff int) Crawler {
	return &ImdbCrawler{
		cutoff: cutoff,
		topUrl: fmt.Sprintf("http://www.imdb.com/search/title?release_date=%v&primary_language=%v&view=simple&ref_=rlm_yr", year, language),
	}
}

func (self *ImdbCrawler) GetBaseUrl() string {
	return self.topUrl
}

func (self *ImdbCrawler) GatherLinks(page *util.Page, ch chan CrawlInfo) {
	page.Document.Find(".lister-col-wrapper").Each(func(i int, lineItem *goquery.Selection) {
		ratingFloat := getRating(lineItem)
		name, link := page.ParseAnchor(lineItem.Find("a"))
		ch <- &ImdbInfo{Name: strings.TrimSuffix(name, "12345678910X"), Link: link, Rating: ratingFloat, CutOff: self.cutoff}
	})
}

func (self *ImdbCrawler) NextPageLink(page *util.Page) (url string, ok bool) {
	var params string
	nextPageElement := page.Document.Find(".next-page")
	if params, ok = nextPageElement.Attr(util.HREF); ok {
		url = self.getImdbUrl(page, params)
	}
	return
}

func (self *ImdbCrawler) PrintSet(good []CrawlInfo, bad []CrawlInfo) bool {
	return true
}

func (self *ImdbCrawler) getImdbUrl(page *util.Page, params string) string {
	return fmt.Sprintf("http://%v%v%v", page.Document.Url.Host, page.Document.Url.Path, params)
}

/* Helpers */
func getRating(lineItem *goquery.Selection) float64 {
	ratingElement := lineItem.Find(".col-imdb-rating > strong")
	rating := strings.TrimSpace(ratingElement.Text())
	if ratingFloat, err := strconv.ParseFloat(rating, 32); err == nil {
		return ratingFloat
	} else {
		log.WithFields(log.Fields{"Error": err}).Error("Error Parsing Rating")
		return -1
	}
}
