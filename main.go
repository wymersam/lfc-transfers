package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func sendToDiscord(webhookURL, title, headline, link string) {
	// Creates a nice-looking card in Discord
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       title,
				"description": headline + "\n\n[Read Article](" + link + ")",
				"color":       15158272, // LFC Red
				"thumbnail": map[string]string{
					"url": "https://upload.wikimedia.org/wikipedia/en/thumb/0/0c/Liverpool_FC.svg/1200px-Liverpool_FC.svg.png",
				},
			},
		},
	}

	jsonData, _ := json.Marshal(payload)
	http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
}

func main() {

	newsSource := []NewsSource{
		{Name: "BBC Sport", Link: "https://www.bbc.com/sport/football/teams/liverpool"},
		{Name: "Liverpool Echo", Link: "https://www.liverpoolecho.co.uk/all-about/liverpool-fc"},
	}

	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	if webhookURL == "" {
		fmt.Println("❌ Error: DISCORD_LFC_WEBHOOK not found in environment!")
		return
	}

	if news := ScrapeLiverpoolNews(newsSource); len(news) > 0 {
		for _, item := range news {
			headline := addColourToHeadline(item.Headline)
			println(item.Source + ": " + headline)

			// TRIGGER POP-UP:
			if strings.Contains(strings.ToLower(item.Headline), "gossip") ||
				strings.Contains(strings.ToLower(item.Headline), "transfer") {

				sendToDiscord(webhookURL, "🚨 LFC Transfer Gossip", item.Headline, item.Link)

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
