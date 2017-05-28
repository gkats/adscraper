package scraper

import (
	"errors"
	"net/http"
)

const UA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36"

type crawler struct{}

func (c *crawler) Fetch(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UA)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, err
	}

	// TODO check other status codes, and redirects
	if res.StatusCode != 200 {
		return res, errors.New("Received status: " + res.Status)
	}
	return res, nil
}
