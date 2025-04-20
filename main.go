package main

import (
	"encoding/json"
	"ingest_bot/app"
	"ingest_bot/helpers"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	bot, err := tgbotapi.NewBotAPI(helpers.Settings.Token)
	if err != nil {
		panic(err)
	}
	//bot.Debug = true
	log.Printf("%s bot has started working", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	upd := bot.GetUpdatesChan(u)
	var followUpMsg *tgbotapi.Message

	for update := range upd {
		updateJSON, err := json.MarshalIndent(update, "", "    ")
		if err != nil {
			log.Printf("failed to marshal update: %v", err)
			continue
		}

		log.Println(string(updateJSON))

		if update.Message == nil && update.CallbackQuery == nil {
			log.Printf("update.Message and CallbackQuery is nil: %v", update.Message)
			continue
		}

		res, err := app.CheckStringMatching(&update)
		msg := update.Message

		if res == true {

			err := app.CallIngestCmd(msg, bot) // forwarding request message
			if err != nil {
				log.Printf("Error calling CallIngestCmd: %s", err)
				helpers.SendMeInfo(err.Error(), bot)
				continue
			}

			followUpMsg, err = app.ConstructCallKeyboard(bot, &update) // sending additional message with keyboard
			if err != nil {
				log.Printf("Error calling ConstructCallKeyboard: %s", err)
				helpers.SendMeInfo(err.Error(), bot)
				continue
			}

		}

		if update.CallbackData() != "" {
			err = app.AnswerCallback(&update, bot, followUpMsg)
			if err != nil {
				continue
			}
		}
	}
}
