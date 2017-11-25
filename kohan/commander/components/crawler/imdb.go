package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"github.com/amanhigh/go-fun/util"
	. "github.com/amanhigh/go-fun/models/crawler"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

type ImdbCrawler struct {
	cutoff int
	topUrl string

	infoChannel chan ImdbInfo
	passedInfos []ImdbInfo
	failedInfos []ImdbInfo
}

func NewImdbCrawler(year int, language string, cutoff int) Crawler {
	return &ImdbCrawler{
		cutoff:      cutoff,
		topUrl:      fmt.Sprintf("http://www.imdb.com/search/title?release_date=%v&primary_language=%v&view=simple&ref_=rlm_yr", year, language),
		infoChannel: make(chan ImdbInfo, 512),
		passedInfos: []ImdbInfo{},
		failedInfos: []ImdbInfo{},
	}
}

func (self *ImdbCrawler) GetBaseUrl() string {
	return self.topUrl
}

func (self *ImdbCrawler) GatherLinks(page *util.Page) {
	page.Document.Find(".lister-col-wrapper").Each(func(i int, lineItem *goquery.Selection) {
		ratingFloat := getRating(lineItem)
		name, link := page.ParseAnchor(lineItem.Find("a"))
		self.infoChannel <- ImdbInfo{Name: strings.TrimSuffix(name, "12345678910X"), Link: link, Rating: ratingFloat}
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

func (self *ImdbCrawler) GatherComplete() {
	close(self.infoChannel)
}

func (self *ImdbCrawler) BuildSet() {
	/* Fire Parallel Consumer to Separate Movies */
	for value := range self.infoChannel {
		if value.Rating >= float64(self.cutoff) || value.Rating < 0.1 {
			self.passedInfos = append(self.passedInfos, value)
		} else {
			self.failedInfos = append(self.failedInfos, value)
		}
	}
}

func (self *ImdbCrawler) PrintSet() {
	/* Output Good/Bad Movies in Separate Sections */
	util.PrintYellow("Passed Info")
	urls := []string{}
	for _, info := range self.passedInfos {
		info.Print()
		urls = append(urls, info.Link)
	}
	ioutil.WriteFile(GOOD_URL_FILE, []byte(strings.Join(urls, "\n")), util.DEFAULT_PERM)

	util.PrintYellow("Failed Info")
	urls = []string{}
	for _, info := range self.failedInfos {
		info.Print()
		urls = append(urls, info.Link)
	}
	ioutil.WriteFile(BAD_URL_FILE, []byte(strings.Join(urls, "\n")), util.DEFAULT_PERM)
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
