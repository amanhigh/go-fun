package util

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/html"
)

const HREF = "href"

type Page struct {
	Document *goquery.Document
}

func NewPageUsingClient(rawUrl string, client HttpClientInterface) (page *Page) {
	response := ""
	if _, err := client.DoGet(rawUrl, &response); err == nil {
		if root, err := html.Parse(strings.NewReader(response)); err == nil {
			doc := goquery.NewDocumentFromNode(root)
			doc.Url, _ = url.Parse(rawUrl)
			page = &Page{}
			page.Document = doc
		} else {
			log.WithFields(log.Fields{"Error": err}).Error("Error Parsing Response")
		}
	} else {
		log.WithFields(log.Fields{"URL": rawUrl, "Error": err}).Error("Error Querying URL")
	}
	return
}

func NewPage(url string) *Page {
	if doc, err := goquery.NewDocument(url); err == nil {
		return &Page{Document: doc}
	} else {
		log.WithFields(log.Fields{"Error": err}).Fatal("Unable to Create Page")
		return nil
	}
}

func (self *Page) ParseAnchor(anchor *goquery.Selection) (text string, link string) {
	var ok bool
	text = anchor.Text()
	if link, ok = anchor.Attr(HREF); ok {
		link = fmt.Sprintf("http://%v%v", self.Document.Url.Host, link)
	}
	return
}
