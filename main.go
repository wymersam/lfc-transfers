package main

import (
	"fmt"
	"net/http"
	"net/url"
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

func sendToTelegram(token, chatID, text string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	// Telegram uses standard URL values
	data := url.Values{
		"chat_id":    {chatID},
		"text":       {text},
		"parse_mode": {"Markdown"}, // Lets you use bold/links
	}

	http.PostForm(apiURL, data)
}
func main() {

	newsSource := []NewsSource{
		{Name: "BBC Sport", Link: "https://www.bbc.com/sport/football/teams/liverpool"},
		{Name: "Liverpool Echo", Link: "https://www.liverpoolecho.co.uk/all-about/liverpool-fc"},
	}

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")

	if telegramToken == "" || telegramChatID == "" {
		fmt.Println("❌ Error: TELEGRAM_TOKEN or TELEGRAM_CHAT_ID not found in environment!")
		return
	}

	if news := ScrapeLiverpoolNews(newsSource); len(news) > 0 {
		for _, item := range news {
			headline := addColourToHeadline(item.Headline)
			println(item.Source + ": " + headline)

			// TRIGGER POP-UP:
			if strings.Contains(strings.ToLower(item.Headline), "gossip") ||
				strings.Contains(strings.ToLower(item.Headline), "transfer") {

				sendToTelegram(telegramToken, telegramChatID, "🚨 LFC Transfer Gossip\n"+item.Headline+"\n\n[Read Article]("+item.Link+")")

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
