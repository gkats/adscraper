package main

import (
	"github.com/gkats/scraper"
)

func main() {
	store, err := scraper.NewStore("dbUrl")
	handleError(err)

	defer Cleanup()
	scraper.NewServer(store).Listen(3000)
}

func Cleanup() {}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
