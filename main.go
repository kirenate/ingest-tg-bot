package main

import (
	"ingest_bot/app"
	"ingest_bot/helpers"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var followUpMsg *tgbotapi.Message
var newUpd tgbotapi.Update

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

	for update := range upd {
		helpers.PrintJSON("", update)

		if update.Message == nil && update.CallbackQuery == nil {
			log.Printf("update.Message and CallbackQuery is nil: %v", update.Message)
			continue
		}

		msg := update.Message
		if msg != nil {
			if msg.Chat.ID == helpers.Settings.FromChatId {
				res, err := app.CheckStringMatching(msg)
				if err != nil {
					log.Printf("String didn't match, %s", err)
				}

				if res {
					app.SendRequestMsgCopy(msg, helpers.Settings.ToChatId, bot) // forwarding request message

					followUpMsg, err = app.ConstructCallKeyboard(bot, &update) // sending additional message with keyboard

					if err != nil {
						log.Printf("Error calling ConstructCallKeyboard: %s", err)
						helpers.SendMeInfo(err.Error(), bot)
					}
				}
			}
		}
		log.Println("before last step")
		helpers.PrintJSON("my update.CallbackData() : ", update.CallbackData())
		if update.CallbackQuery != nil {
			newUpd.CallbackQuery = update.CallbackQuery
		}
		helpers.PrintJSON("my newUpd: ", newUpd)
		if newUpd.CallbackData() != "" {
			log.Println("inside last step")
			err = app.AnswerCallback(&newUpd, bot, followUpMsg)
			if err != nil {
				continue
			}
			log.Println("after last step")
		}
	}
}
