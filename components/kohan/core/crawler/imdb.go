package crawler

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	clients2 "github.com/amanhigh/go-fun/common/clients"
	util2 "github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/crawler"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
)

type ImdbCrawler struct {
	cutoff   int
	language string
	topUrl   string
	client   *resty.Client
}

func NewImdbCrawler(year int, language string, cutoff int, cookies string) Crawler {
	color.Yellow("ImdbCrawler: Year:%v Lang:%v Cutoff: %v", year, language, cutoff)

	if util2.IsDebugMode() {
		fmt.Println("IMDB Cookie: ", cookies)
	}
	//TODO: enable Compression
	client := clients2.NewHttpClientWithCookies("https://www.imdb.com", util2.ParseCookies(cookies), clients2.DefaultHttpClient)
	return &ImdbCrawler{
		cutoff:   cutoff,
		language: language,
		topUrl:   fmt.Sprintf("https://www.imdb.com/search/title?release_date=%v&primary_language=%v&view=simple&title_type=feature&sort=num_votes,desc", year, language),
		client:   client,
	}
}

func (self *ImdbCrawler) GetTopPage() *util2.Page {
	topPage := util2.NewPageUsingClient(self.topUrl, self.client)
	userName := topPage.Document.Find("a.navbar__user-name").Text()
	if strings.Contains(userName, "Aman") {
		color.Green("User: %s", userName)
	} else {
		color.Red("User Not Logged in Please Check Cookie is Present and Not Stale")
	}
	return topPage
}

func (self *ImdbCrawler) GatherLinks(page *util2.Page, ch chan crawler.CrawlInfo) {
	page.Document.Find(".lister-col-wrapper").Each(func(i int, lineItem *goquery.Selection) {
		/* Read Rating & Link from List Page */
		ratingFloat := getRating(lineItem)
		name, link := page.ParseAnchor(lineItem.Find(".lister-item-header a"))

		/* Go Crawl Movie Page for My Rating & Other Details */
		if moviePage := util2.NewPageUsingClient(link, self.client); moviePage != nil {
			myRating := util2.ParseFloat(moviePage.Document.Find(".star-rating-value").Text())

			ch <- &crawler.ImdbInfo{
				Name: name,
				Link: link, Rating: ratingFloat,
				Language: self.language,
				MyRating: myRating,
				CutOff:   self.cutoff,
			}
		}
	})
}

func (self *ImdbCrawler) NextPageLink(page *util2.Page) (url string, ok bool) {
	var params string
	nextPageElement := page.Document.Find(".next-page")
	if params, ok = nextPageElement.Attr(util2.HREF); ok {
		url = fmt.Sprintf("https://%v%v", page.Document.Url.Host, params)
	}
	return
}

func (self *ImdbCrawler) PrintSet(good CrawlSet, bad CrawlSet) bool {
	return true
}

func (self *ImdbCrawler) getImdbUrl(page *util2.Page, params string) string {
	return fmt.Sprintf("https://%v%v%v", page.Document.Url.Host, page.Document.Url.Path, params)
}

/* Helpers */
func getRating(lineItem *goquery.Selection) float64 {
	ratingElement := lineItem.Find(".col-imdb-rating > strong")
	return util2.ParseFloat(ratingElement.Text())
}
