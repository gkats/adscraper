package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gkats/adscraper"
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

	client := adscraper.NewClient(hostUrl)
	ks, err := client.GetKeywords()
	handleError(err)

	// Scrape ads for each keyword
	for _, k := range ks {
		ads, err := adscraper.Scrape(adscraper.NewURL(k.Value))
		handleError(err)

		// POST each ad to the ads service
		for _, ad := range ads {
			handleError(client.PostAdKeywords(ad, k))
		}
		// PATCH to increment keyword scraped attributes
		handleError(client.PatchKeyword(k.ID))
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
