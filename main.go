package main

import (
	"fmt"
	"strings"

	"github.com/gen2brain/beeep"
)

func addColourToHeadline(headline string) string {
	if strings.Contains(strings.ToLower(headline), "transfer") {
		return "\033[32m" + headline + "\033[0m" // Green text for transfer news
	}

	if strings.Contains(strings.ToLower(headline), "liverpool v") {
		return "\033[34m" + headline + "\033[0m" // Purple text for Liverpool match news
	}

	return headline
}

func main() {

	newsSource := []NewsSource{
		{Name: "BBC Sport", Link: "https://www.bbc.com/sport/football/teams/liverpool"},
		{Name: "Liverpool Echo", Link: "https://www.liverpoolecho.co.uk/all-about/liverpool-fc"},
	}

	if news := ScrapeLiverpoolNews(newsSource); len(news) > 0 {
		for _, item := range news {
			headline := addColourToHeadline(item.Headline)
			println(item.Source + ": " + headline)

			// TRIGGER POP-UP:
			if strings.Contains(strings.ToLower(item.Headline), "gossip") ||
				strings.Contains(strings.ToLower(item.Headline), "transfer") {

				err := beeep.Notify(
					"LFC Scout: "+item.Source, // Title
					item.Headline,             // Message body
					"/Users/sammy-jo.wymer@diconium.com/Documents/dev/personal-development/LFC-transfers/images/lfc-icon.png")
				if err != nil {
					fmt.Println("Could not send notification:", err)
				}
			}
		}
	} else {
		println("No LFC news found.")
	}
}
