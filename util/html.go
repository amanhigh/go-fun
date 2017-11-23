package util

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"fmt"
)

const HREF = "href"

type Page struct {
	Document *goquery.Document
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
