package command

import (
	"github.com/amanhigh/go-fun/apps/components/kohan/core/crawler"
	"github.com/amanhigh/go-fun/util"
	"github.com/spf13/cobra"
)

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Different Crawler Commands",
	Args:  cobra.ExactArgs(1),
}

var imdbCmd = &cobra.Command{
	Use:   "imdb [Year] [Language]",
	Short: "Imdb Crawler",
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if year, err = util.ParseInt(args[0]); err == nil {
			util.ValidateEnumArg(args[1], []string{"pa", "en", "hi"})
		}
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		crawler.NewCrawlerManager(crawler.NewImdbCrawler(year, args[1], cutOff, keyFilePath), count, verbose).Crawl()
	},
}

var gameCmd = &cobra.Command{
	Use:   "game [toplink]",
	Short: "Gamespot Crawler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		crawler.NewCrawlerManager(crawler.NewGameSpotCrawler(args[0]), count, verbose).Crawl()
	},
}

var hubCmd = &cobra.Command{
	Use:   "hub [toplink]",
	Short: "Hub Crawler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		crawler.NewCrawlerManager(crawler.NewHubCrawler(args[0]), count, verbose).Crawl()
	},
}

func init() {
	crawlCmd.PersistentFlags().IntVarP(&count, "count", "c", 200, "Count of entries to be crawled")
	crawlCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable Verbose Mode")

	imdbCmd.Flags().IntVarP(&cutOff, "cutoff", "o", 5, "Cut Off For Movie")
	imdbCmd.Flags().StringVarP(&keyFilePath, "path", "p", "/tmp/imdb.yml", "IMDB Yaml Cookie File (id,sid)")

	RootCmd.AddCommand(crawlCmd)
	crawlCmd.AddCommand(imdbCmd, gameCmd, hubCmd)
}
