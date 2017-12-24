package processor

import (
	"flag"
	. "github.com/amanhigh/go-fun/kohan/commander/components/crawler"
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
	requiredCount := flagSet.Int("r", 500, "Required Count")
	verbose := flagSet.Bool("v", false, "Verbose")
	langCode := flagSet.String("l", "en", "Language Code [pa,en,hi]")
	keyFile := flagSet.String("k", "/tmp/imdb.key", "IMDB Key File")
	e := flagSet.Parse(args)
	NewCrawlerManager(NewImdbCrawler(*year, *langCode, *cutoff, *keyFile), *requiredCount, *verbose).Crawl()
	return e
}