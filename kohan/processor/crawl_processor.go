package processor

import (
	"flag"
	"github.com/amanhigh/go-fun/kohan/commander/components/crawler"
)

type CrawlProcessor struct {
}

func (self *CrawlProcessor) GetArgedHandlers() (map[string]HandleFunc) {
	return map[string]HandleFunc{
		"imdb": self.handleImdb,
	}
}

func (self *CrawlProcessor) GetNonArgedHandlers() (map[string]DirectFunc) {
	return map[string]DirectFunc{}
}

func (self *CrawlProcessor) handleImdb(flagSet *flag.FlagSet, args []string) error {
	year := flagSet.Int("y", 2015, "Year of Movie")
	cutoff := flagSet.Int("c", 5, "Movie Cutoff")
	langCode := flagSet.String("l", "en", "Language Code [pa,en,hi]")
	e := flagSet.Parse(args)
	crawl(crawler.NewImdbCrawler(*year, *langCode, *cutoff))
	return e
}

func crawl(crawler crawler.Crawler) {
	crawler.CrawlerManager{Crawler: crawler}.Crawl()
}
