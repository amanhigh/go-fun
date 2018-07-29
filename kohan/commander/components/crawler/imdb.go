package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	. "github.com/amanhigh/go-fun/models/crawler"
	"github.com/amanhigh/go-fun/util"
)

type ImdbCrawler struct {
	cutoff   int
	language string
	topUrl   string
	client   util.HttpClientInterface
}

func NewImdbCrawler(year int, language string, cutoff int, keyFile string) Crawler {
	util.PrintYellow(fmt.Sprintf("ImdbCrawler: Year:%v Lang:%v Cutoff: %v", year, language, cutoff))

	if key, _ := ioutil.ReadFile(keyFile); len(key) > 0 {
		cookie := http.Cookie{Name: "id", Value: string(key)}
		client := util.NewHttpClientWithCookies("http://www.imdb.com", []*http.Cookie{&cookie}, true, true)
		return &ImdbCrawler{
			cutoff:   cutoff,
			language: language,
			topUrl:   fmt.Sprintf("http://www.imdb.com/search/title?release_date=%v&primary_language=%v&view=simple&title_type=feature&sort=num_votes,desc", year, language),
			client:   client,
		}
	} else {
		log.WithFields(log.Fields{"KeyFile": keyFile}).Fatal("Empty IMDB Key")
		return nil
	}
}

func (self *ImdbCrawler) GetTopPage() *util.Page {
	return util.NewPageUsingClient(self.topUrl, self.client)
}

func (self *ImdbCrawler) GatherLinks(page *util.Page, ch chan CrawlInfo) {
	page.Document.Find(".lister-col-wrapper").Each(func(i int, lineItem *goquery.Selection) {
		/* Read Rating & Link from List Page */
		ratingFloat := getRating(lineItem)
		name, link := page.ParseAnchor(lineItem.Find("a"))

		/* Go Crawl Movie Page for My Rating & Other Details */
		if moviePage := util.NewPageUsingClient(link, self.client); moviePage != nil {
			myRating := util.ParseFloat(moviePage.Document.Find(".star-rating-value").Text())

			ch <- &ImdbInfo{
				Name: strings.TrimSuffix(name, "12345678910X"),
				Link: link, Rating: ratingFloat,
				Language: self.language,
				MyRating: myRating,
				CutOff:   self.cutoff,
			}
		}
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
	return util.ParseFloat(ratingElement.Text())
}
