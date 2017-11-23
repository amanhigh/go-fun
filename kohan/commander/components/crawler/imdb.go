package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"github.com/amanhigh/go-fun/util"
	. "github.com/amanhigh/go-fun/models/crawler"
	"sync"
	log "github.com/Sirupsen/logrus"
)

type ImdbCrawler struct {
	cutoff int
	page   *util.Page
}

func NewImdbCrawler(year int, language string, cutoff int) *ImdbCrawler {
	url := fmt.Sprintf("http://www.imdb.com/search/title?release_date=%v&primary_language=%v&view=simple&ref_=rlm_yr", year, language)
	return &ImdbCrawler{
		cutoff: cutoff,
		page:   util.NewPage(url),
	}
}

func (self *ImdbCrawler) Crawl() {
	imdbInfoChannel := make(chan ImdbInfo, 512)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go self.crawlRecursive(self.page, imdbInfoChannel, waitGroup)
	passedInfos := []ImdbInfo{}
	failedInfos := []ImdbInfo{}
	go func() {
		for value := range imdbInfoChannel {
			if value.Rating > float64(self.cutoff) {
				passedInfos = append(passedInfos, value)
			} else {
				failedInfos = append(failedInfos, value)
			}
			value.Print()
		}
	}()
	waitGroup.Wait()
	close(imdbInfoChannel)

	util.PrintYellow("Passed Info")
	for _, info := range passedInfos {
		info.Print()
	}

	util.PrintYellow("Failed Info")
	for _, info := range failedInfos {
		info.Print()
	}
}

func (self *ImdbCrawler) crawlRecursive(page *util.Page, infos chan ImdbInfo, waitGroup sync.WaitGroup) {
	util.PrintYellow("Processing: " + page.Document.Url.String())

	/* If Next Link is Present Crawl It */
	nextPageElement := page.Document.Find(".next-page")
	if params, ok := nextPageElement.Attr(util.HREF); ok {
		nextUrl := self.getImdbUrl(params)
		waitGroup.Add(1)
		go self.crawlRecursive(util.NewPage(nextUrl), infos, waitGroup)
	}

	/* Find Links for this Page */
	self.findLinks(self.page, infos)
	waitGroup.Done()
}

func (self *ImdbCrawler) getImdbUrl(params string) string {
	return fmt.Sprintf("http://%v%v%v", self.page.Document.Url.Host, self.page.Document.Url.Path, params)
}

func (self *ImdbCrawler) findLinks(page *util.Page, infoChannel chan ImdbInfo) {
	page.Document.Find(".lister-col-wrapper").Each(func(i int, lineItem *goquery.Selection) {
		ratingFloat := getRating(lineItem)
		name, link := page.ParseAnchor(lineItem.Find("a"))
		infoChannel <- ImdbInfo{Name: strings.TrimSuffix(name, "12345678910X"), Link: link, Rating: ratingFloat}
	})
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
