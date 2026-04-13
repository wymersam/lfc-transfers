package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

type NewsItem struct {
	Source   string
	Headline string
	Link     string
	Time     string
}

type NewsSource struct {
	Name string
	Link string
}

func ScrapeLiverpoolNews(newsSources []NewsSource) []NewsItem {
	var news []NewsItem
	var mu sync.Mutex
	var wg sync.WaitGroup

	// We create a collector inside the function or use one per goroutine
	// to ensure thread safety with Colly
	for _, newsSource := range newsSources {
		wg.Add(1)

		go func(src NewsSource) {
			defer wg.Done()

			c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Go-Scout)"))

			c.OnError(func(r *colly.Response, err error) {
				fmt.Printf("Error visiting %s: %v\n", r.Request.URL, err)
			})

			// BBC uses h3 inside article cards; also try broader selectors
			c.OnHTML("h3, h2, [data-testid='card-headline']", func(e *colly.HTMLElement) {
				headline := strings.TrimSpace(e.Text)
				link := e.Attr("href")

				// Filter for Liverpool keywords
				if strings.Contains(strings.ToLower(headline), "liverpool") ||
					strings.Contains(strings.ToLower(headline), "lfc") ||
					strings.Contains(strings.ToLower(headline), "transfer") {

					item := NewsItem{
						Source:   e.Request.URL.Host,
						Headline: strings.TrimSpace(headline),
						Link:     e.Request.AbsoluteURL(link),
					}
					mu.Lock()
					news = append(news, item)
					mu.Unlock()
				}
			})
			c.Visit(src.Link)
		}(newsSource)
	}

	wg.Wait()
	return news
}
