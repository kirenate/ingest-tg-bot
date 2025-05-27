package app

import (
	"errors"
	"fmt"
	"ingest_bot/helpers"
	"regexp"

	"github.com/rs/zerolog/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ConstructCallKeyboard(bot *tgbotapi.BotAPI, update *tgbotapi.Update) (*tgbotapi.Message, error) {
	follow_up_msg := tgbotapi.NewMessage(helpers.Settings.ToChatId, "Нас вызывают!")
	follow_up_msg.ReplyMarkup = AcceptKeyboard
	resp, err := bot.Send(follow_up_msg)
	if err != nil {
		log.Error().Stack().Err(err).Msg("'Send' could not be executed")
	}
	return &resp, err

}

func CheckStringMatching(msg *tgbotapi.Message) (bool, error) {
	if msg == nil {
		err := errors.New("update.Message is nil")
		return false, err
	}
	r, err := regexp.Compile(`.*(?i)прокс.*|.*(?i)инжест.*|.*(?i)proxy.*`)
	if err != nil {
		log.Panic().Msgf("Regex could not be compiled!\n%s", err)
		return false, err
	}
	res := r.MatchString(msg.Text)
	log.Info().Msgf("String matching: %v", res)
	return res, err
}

func GetCallbackQueryResponse(update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	if update.CallbackQuery == nil {
		err := errors.New("update.CallbackQuery is nil")
		return err
	}
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	_, err := bot.Request(callback)
	//log.Printf("Callback msg response: %v", resp)
	if err != nil {
		log.Printf("Error while getting accepting callback data: %s", err)
		//helpers.SendMeInfo(err.Error(), bot)
		return err
	}
	return err
}

func EditFollowUpMessage(respChatID int64, respMessageID int, bot *tgbotapi.BotAPI) (*tgbotapi.Message, error) {
	editedKeyboard := tgbotapi.NewEditMessageReplyMarkup(respChatID, respMessageID, ConfirmKeyboard)
	newResp, err := bot.Send(editedKeyboard)
	if err != nil {
		log.Printf("Error while sending editedKeyboard to telegram : %s", err)
		return nil, err
	}

	return &newResp, err
}

func SendMsgConfirmation(respChatID int64, respMessageID int, userId int64, bot *tgbotapi.BotAPI) (*tgbotapi.Message, error) {
	customMessage := helpers.ChooseCustomMessage(userId)
	msg_confirmed := tgbotapi.NewEditMessageText(respChatID, respMessageID, customMessage)
	new_msg, err := bot.Send(msg_confirmed)
	if err != nil {
		log.Printf("Error sending confirmation message: %s", err)
		return nil, err
	}
	return &new_msg, err
}

func SendRequestMsgCopy(requestMsg *tgbotapi.Message, ToChatId int64, bot *tgbotapi.BotAPI) string {
	requestCopy := fmt.Sprintf("@%s\n\n%s", requestMsg.From.UserName, requestMsg.Text)
	msgCopy := tgbotapi.NewMessage(helpers.Settings.ToChatId, requestCopy)
	_, err := bot.Send(msgCopy)
	if err != nil {
		helpers.SendMeInfo(err.Error(), bot)
		log.Printf("Message could not be copied, %s", err)
	}
	return requestCopy
}

func AnswerCallback(newUpd *tgbotapi.Update, bot *tgbotapi.BotAPI, followUpMsg *tgbotapi.Message) error {
	if newUpd.CallbackQuery == nil {
		err := errors.New("update.CallbackQuery is nil")
		return err
	}
	switch newUpd.CallbackQuery.Data {
	case "request_accepted":
		err := GetCallbackQueryResponse(newUpd, bot)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
		}
		if newUpd == nil {
			log.Panic().Msg("newUpd is nil")
		}
		if newUpd.CallbackQuery == nil {
			log.Panic().Msg("newUpd.Message is nil")
		}
		_, err = EditFollowUpMessage(newUpd.CallbackQuery.Message.Chat.ID, newUpd.CallbackQuery.Message.MessageID, bot)

		if err != nil {
			log.Error().Stack().Err(err).Msg("error calling EditFollowUpMessage")
		}
	case "request_satisfied":
		err := GetCallbackQueryResponse(newUpd, bot)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Error calling GetCallbackQueryResponse")
		}
		if newUpd == nil {
			log.Panic().Msg("newUpd is nil")
		}
		if newUpd.CallbackQuery == nil {
			log.Panic().Msg("newUpd.Message is nil")
		}
		new_msg, err := SendMsgConfirmation(newUpd.CallbackQuery.Message.Chat.ID, newUpd.CallbackQuery.Message.MessageID, newUpd.CallbackQuery.From.ID, bot)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error sending confirmation")
		}
		log.Info().Interface("msg", new_msg).Msg("telegram returned message:")
	default:
		log.Info().Msg("No valid callback data found")
	}
	return nil
}

//////////////////////////////////////////////

var AcceptKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Я иду помогать!", "request_accepted"),
	),
)

var ConfirmKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Миссия выполнена", "request_satisfied"),
	),
)
