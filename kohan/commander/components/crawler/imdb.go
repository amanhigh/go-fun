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
	"io/ioutil"
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
	/* Build Channel & WG for Parallel Parsers */
	imdbInfoChannel := make(chan ImdbInfo, 512)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	/* Fire First Crawler */
	go self.crawlRecursive(self.page, imdbInfoChannel, waitGroup)

	/* Fire Parallel Consumer to Separate Movies */
	var passedInfos []ImdbInfo
	var failedInfos []ImdbInfo
	go func() {
		for value := range imdbInfoChannel {
			if value.Rating > float64(self.cutoff) || value.Rating < 0.1 {
				passedInfos = append(passedInfos, value)
			} else {
				failedInfos = append(failedInfos, value)
			}
		}
	}()

	/* Wait Till all Parsers Complete & Close Channel */
	waitGroup.Wait()
	close(imdbInfoChannel)

	/* Output Good/Bad Movies in Separete Sections */
	util.PrintYellow("Passed Info")
	urls := []string{}
	for _, info := range passedInfos {
		info.Print()
		urls = append(urls, info.Link)
	}
	ioutil.WriteFile(GOOD_URL_FILE, []byte(strings.Join(urls, "\n")), util.DEFAULT_PERM)

	util.PrintYellow("Failed Info")
	urls = []string{}
	for _, info := range failedInfos {
		info.Print()
		urls = append(urls, info.Link)
	}
	ioutil.WriteFile(BAD_URL_FILE, []byte(strings.Join(urls, "\n")), util.DEFAULT_PERM)
}

/**
	Recursively Crawl Given Page moving to next if next link is available.
	Write all Movies of current page onto channel
 */
func (self *ImdbCrawler) crawlRecursive(page *util.Page, infos chan ImdbInfo, waitGroup *sync.WaitGroup) {
	util.PrintYellow("Processing: " + page.Document.Url.String())

	/* If Next Link is Present Crawl It */
	nextPageElement := page.Document.Find(".next-page")
	if params, ok := nextPageElement.Attr(util.HREF); ok {
		nextUrl := self.getImdbUrl(params)
		waitGroup.Add(1)
		go self.crawlRecursive(util.NewPage(nextUrl), infos, waitGroup)
	}

	/* Find Links for this Page */
	self.findLinks(page, infos)
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
