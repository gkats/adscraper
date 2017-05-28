package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

const GOOGLE_URL = "https://www.google.com/search"

func NewURL(s string) string {
	return GOOGLE_URL + "?q=" + strings.Replace(s, " ", "+", -1)
}

func Scrape(url string) ([]*Ad, error) {
	c := &crawler{}

	res, err := c.Fetch(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return extract(res)
}

func extract(r *http.Response) ([]*Ad, error) {
	var ads = make([]*Ad, 0)

	doc, err := goquery.NewDocumentFromResponse(r)
	if err != nil {
		return ads, err
	}
	doc.Find(".ads-ad").Each(func(i int, sel *goquery.Selection) {
		ads = append(ads, extractAd(i+1, sel))
	})

	return ads, err
}

func extractAd(pos int, sel *goquery.Selection) *Ad {
	ad := &Ad{}

	ad.H1, ad.H2 = splitHead(sel.Find("h3 > a").Text())
	ad.Path = strings.TrimSpace(sel.Find(".ads-visurl cite").Text())
	ad.Desc = strings.TrimSpace(sel.Find(".ads-creative").Text())
	ad.SetRest(innerHTML(sel.Find(".ads-creative")))
	ad.SetRaw(innerHTML(sel))
	ad.Position = pos

	return ad
}

func splitHead(head string) (string, string) {
	h := strings.SplitN(head, "-", 2)
	if len(h) > 1 {
		return normalize(strings.TrimSpace(h[0])), normalize(strings.TrimSpace(h[1]))
	} else {
		return "", ""
	}
}

func normalize(s string) string {
	return strings.Replace(s, "\u200e", "", -1)
}

func innerHTML(sel *goquery.Selection) string {
	html := make([]string, 0)
	sel.NextAll().Each(func(i int, sel *goquery.Selection) {
		h, _ := goquery.OuterHtml(sel)
		html = append(html, strings.TrimSpace(h))
	})
	return strings.Join(html, "")
}
