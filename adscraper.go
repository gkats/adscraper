package adscraper

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func NewURL(s string) string {
	return "https://www.google.com/search?q=" + strings.Replace(s, " ", "+", -1)
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

	ad.Position = pos
	ad.H1, ad.H2 = splitHead(sel.Find("h3 > a").Text())
	ad.Path = strings.TrimSpace(sel.Find(".ads-visurl cite").Text())

	descSel := sel.Find(".ads-creative")
	ad.Desc = strings.TrimSpace(descSel.Text())
	ad.SetRest(innerHTML(descSel))

	raw, _ := goquery.OuterHtml(sel)
	ad.SetRaw(raw)

	return ad
}

func splitHead(head string) (h1 string, h2 string) {
	h1, h2 = "", ""
	if h := strings.SplitN(head, "-", 2); len(h) > 1 {
		h1 = normalize(strings.TrimSpace(h[0]))
		h2 = normalize(strings.TrimSpace(h[1]))
	}
	return
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
