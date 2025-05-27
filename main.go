package main

import (
	"errors"
	"ingest_bot/app"
	"ingest_bot/helpers"
	"ingest_bot/logger"

	"github.com/rs/zerolog/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var followUpMsg *tgbotapi.Message
var newUpd tgbotapi.Update

func main() {
	logger.MakeLogger()
	bot, err := tgbotapi.NewBotAPI(helpers.Settings.Token)
	if err != nil {
		panic(err)
	}
	//bot.Debug = true
	log.Info().Msgf("%s bot has started working", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	upd := bot.GetUpdatesChan(u)

	for update := range upd {
		log.Info().Interface("upd", upd).Msg("update.received")
		//helpers.PrintJSON("", update)

		if update.Message == nil && update.CallbackQuery == nil {
			err = errors.New("update.Message AND CallbackQuery is nil")
			log.Error().Stack().Err(err).Msg("")
			continue
		}

		msg := update.Message
		if msg != nil {
			processMsg(msg, bot, update)
		}
		if update.CallbackQuery != nil {
			newUpd.CallbackQuery = update.CallbackQuery
			log.Info().Interface("CallbackQuery", newUpd.CallbackQuery).Msg("")
			err = app.AnswerCallback(&newUpd, bot, followUpMsg)
			if err != nil {
				log.Info().Msgf("%s", err)
				continue
			}
		}
	}
}

func processMsg(msg *tgbotapi.Message, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if msg.Chat.ID == helpers.Settings.FromChatId {
		res, err := app.CheckStringMatching(msg)
		if err != nil {
			log.Error().Stack().Err(err).Msg("String didn't match")
		}

		if res {
			copy := app.SendRequestMsgCopy(msg, helpers.Settings.ToChatId, bot)
			log.Info().Msgf("message sent:\n%s", copy)

			followUpMsg, err = app.ConstructCallKeyboard(bot, &update) // sending additional message with keyboard
			if err != nil {
				log.Error().Stack().Err(err).Msg("")
				helpers.SendMeInfo(err.Error(), bot)
			} else {
				log.Info().Msg("keyboard attached")
			}
		}
	}
}
