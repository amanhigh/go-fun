package util

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const HREF = "href"

type Page struct {
	Document *goquery.Document
}

func NewPageUsingClient(rawUrl string, client *resty.Client) (page *Page) {
	response := ""
	if _, err := client.R().SetResult(&response).Get(rawUrl); err == nil {
		if root, err := html.Parse(strings.NewReader(response)); err == nil {
			doc := goquery.NewDocumentFromNode(root)
			doc.Url, _ = url.Parse(rawUrl)
			page = &Page{}
			page.Document = doc
		} else {
			log.Error().Err(err).Msg("Error Parsing Response")
		}
	} else {
		log.Error().Str("URL", rawUrl).Err(err).Msg("Error Querying URL")
	}
	return
}

func NewPage(url string) *Page {
	if doc, err := goquery.NewDocument(url); err == nil {
		return &Page{Document: doc}
	} else {
		log.Fatal().Err(err).Msg("Unable to Create Page")
		return nil
	}
}

func (self *Page) ParseAnchor(anchor *goquery.Selection) (text string, link string) {
	var ok bool
	text = anchor.Text()
	if link, ok = anchor.Attr(HREF); ok {
		link = fmt.Sprintf("https://%v%v", self.Document.Url.Host, link)
	}
	return
}

func ParseCookies(rawCookies string) (cookies []*http.Cookie) {
	header := http.Header{}
	header.Add("Cookie", rawCookies)
	request := http.Request{Header: header}
	cookies = request.Cookies()
	return
}
