package main

import (
	"flag"
	"fmt"
	"github.com/gkats/scraper"
	"os"
)

func main() {
	var (
		hostUrl string
	)
	flag.StringVar(&hostUrl, "h", "", "Base URL for the ads service host.")
	flag.Parse()
	if hostUrl == "" {
		fmt.Fprintf(os.Stderr, "You must provide the ads service host URL. Run with --help to see usage instructions.\n")
		os.Exit(1)
	}

	// TODO Get 20 random keywords
	// k, err := store.Keywords.Get()
	// handleError(err)
	keywords := make([]*scraper.Keyword, 1)
	keywords[0] = scraper.New("new reebok shoes")

	client := scraper.NewClient(hostUrl)

	// Scrape ads for each keyword
	for _, k := range keywords {
		ads, err := scraper.Scrape(scraper.NewURL(k.Value))
		handleError(err)

		// POST each ad to the ads service
		for _, ad := range ads {
			handleError(client.PostAdKeywords(ad, k))
		}
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
