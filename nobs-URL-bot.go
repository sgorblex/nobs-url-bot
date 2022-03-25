package main

import (
	"bufio"
	nobs "github.com/sgorblex/nobs-url/lib"
	"log"
	"os"

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
		if update.InlineQuery == nil { // if no inline query, ignore it
			continue
		}

		repl, ok := nobs.Cleanup(update.InlineQuery.Query)
		if ok {
			article := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, "Remove bs", repl)
			article.Description = update.InlineQuery.Query

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
