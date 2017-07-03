package main

import (
	"flag"
	"github.com/gkats/scraper"
)

func main() {
	var (
		dbUrl string
	)
	flag.StringVar(&dbUrl, "d", "", "The database URL. Should be in 'user:password@host:port/database' format.")
	flag.Parse()

	store, err := scraper.NewStore(dbUrl)
	handleError(err)
	defer store.Close()

	scraper.NewServer(store).Listen(3000)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
