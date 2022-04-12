package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	nobs "github.com/sgorblex/nobs-url/lib"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getApiKey() string {
	keyFile, err := os.Open("api_key.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer keyFile.Close()

	sc := bufio.NewScanner(keyFile)
	if !sc.Scan() {
		log.Panic("ERROR: API Key file is empty.")
		os.Exit(1)
	}
	return sc.Text()
}

func main() {
	bot, err := tgbotapi.NewBotAPI(getApiKey())
	if err != nil {
		log.Panic(err)
	}
	// bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.InlineQuery != nil {
			query := strings.SplitN(strings.TrimSpace(update.InlineQuery.Query), " ", 2)
			fmt.Println(query)
			if len(query) < 1 || !nobs.IsURL(query[0]) {
				continue
			}
			clean, ok := nobs.Cleanup(query[0])
			if ok {
				var repl, desc string
				if len(query) == 2 {
					repl = "[" + query[1] + "](" + clean + ")"
					desc = query[1] + ": " + clean
				} else {
					repl = clean
					desc = clean
				}
				article := tgbotapi.NewInlineQueryResultArticleMarkdown(update.InlineQuery.ID, "Remove bs", repl)
				article.Description = desc

				inlineConf := tgbotapi.InlineConfig{
					InlineQueryID: update.InlineQuery.ID,
					IsPersonal:    true,
					CacheTime:     0,
					Results:       []interface{}{article},
				}

				if _, err := bot.Request(inlineConf); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
