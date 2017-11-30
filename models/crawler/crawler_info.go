package crawler

type CrawlInfo interface {
	GoodBad() bool
	ToUrl() string
	Print()
}
