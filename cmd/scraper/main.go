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
		dbUrl   string
	)
	flag.StringVar(&hostUrl, "h", "", "Base URL for the ads service host.")
	flag.StringVar(&dbUrl, "d", "", "The PostgreSQL database URL. Should be in 'user:password@host:port/database' format.")
	flag.Parse()
	if hostUrl == "" {
		fmt.Fprintf(os.Stderr, "You must provide the ads service host URL. Run with --help to see usage instructions.\n")
		os.Exit(1)
	}

	store, err := scraper.NewStore(dbUrl)
	handleError(err)
	defer store.Close()

	client := scraper.NewClient(hostUrl)
	ks, err := client.GetKeywords()
	handleError(err)

	// Scrape ads for each keyword
	for _, k := range ks {
		ads, err := scraper.Scrape(scraper.NewURL(k.Value))
		handleError(err)

		// POST each ad to the ads service
		for _, ad := range ads {
			handleError(client.PostAdKeywords(ad, k))
		}
		// PATCH to increment keyword scraped attributes
		handleError(client.PatchKeyword(k.Id))
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
