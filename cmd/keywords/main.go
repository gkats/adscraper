package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gkats/scraper/keywords"
	"os"
)

func main() {
	var (
		filename string
		dbUrl    string
	)
	flag.StringVar(&filename, "f", "", "Absolute path to the keywords file.")
	flag.StringVar(&dbUrl, "d", "", "The PostgreSQL database URL. Should be in 'user:password@host:port/database' format.")
	flag.Parse()
	if filename == "" {
		fmt.Fprintf(os.Stderr, "You must provide a filename. Run with --help to see usage instructions.\n")
		os.Exit(1)
	}

	// Set up the store
	ks, err := keywords.NewStore(dbUrl)
	handleError(err)
	defer Cleanup(ks)

	// Set up the writer service
	kw := keywords.NewWriter(ks)

	// Open the keywords file
	f, err := os.Open(filename)
	handleError(err)

	// Read file line by line and store each keyword
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		_, err := kw.Upsert(keywords.New(scanner.Text()))
		handleError(err)
	}
	handleError(scanner.Err())
}

func Cleanup(s keywords.Store) error {
	return s.Close()
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
