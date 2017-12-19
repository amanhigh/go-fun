package crawler

type CrawlInfo interface {
	GoodBad() error
	ToUrl() string
	Print()
}
