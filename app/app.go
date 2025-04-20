package app

import (
	"errors"
	"ingest_bot/helpers"
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CallIngestCmd(request_msg *tgbotapi.Message, bot *tgbotapi.BotAPI) error {

	fwd_msg := tgbotapi.NewForward(helpers.Settings.IngestChatId, helpers.Settings.NozhiChatId, request_msg.MessageID)
	//log.Println(fwd_msg)
	_, err := bot.Send(fwd_msg)
	if err != nil {
		log.Printf("Request message could not be forwarded, error: %s", err)
		return err
	}
	return err
}

func ConstructCallKeyboard(bot *tgbotapi.BotAPI, update *tgbotapi.Update) (*tgbotapi.Message, error) {
	follow_up_msg := tgbotapi.NewMessage(helpers.Settings.IngestChatId, "Нас вызывают!")
	follow_up_msg.ReplyMarkup = AcceptKeyboard
	resp, err := bot.Send(follow_up_msg)
	if err != nil {
		log.Printf("'Send' could not be executed: %s", err)
	}

	//log.Printf("my follow up message:\n%v", resp)
	return &resp, err

}

func CheckStringMatching(update *tgbotapi.Update) (bool, error) {
	if update.Message == nil {
		err := errors.New("update.Message is nil")
		return false, err
	}
	r, err := regexp.Compile(`.*(?i)инжест.*`)
	if err != nil {
		log.Panicf("Regex could not be compiled!\n%s", err)
		return false, err
	}
	res := r.MatchString(update.Message.Text)
	log.Printf("String matching: %v", res)
	return res, err
}

func GetCallbackQueryResponse(update *tgbotapi.Update, bot *tgbotapi.BotAPI) (*tgbotapi.APIResponse, error) {
	if update.CallbackQuery == nil {
		err := errors.New("update.CallbackQuery is nil")
		return nil, err
	}
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	resp, err := bot.Request(callback)

	if err != nil {
		log.Printf("Error while getting accepting callback data: %s", err)
		helpers.SendMeInfo(err.Error(), bot)
		return nil, err
	}
	return resp, err
}

func EditFollowUpMessage(followUpMsg *tgbotapi.Message, bot *tgbotapi.BotAPI) (*tgbotapi.Message, error) {
	if followUpMsg == nil {
		log.Panic("followUpMsg is nil")

	}
	if followUpMsg.Chat == nil {
		log.Panic("followUpMsg.Chat is nil")
	}
	editedKeyboard := tgbotapi.NewEditMessageReplyMarkup(followUpMsg.Chat.ID, followUpMsg.MessageID, ConfirmKeyboard)
	resp, err := bot.Send(editedKeyboard)
	if err != nil {
		log.Printf("Error while sending editedKeyboard to telegram : %s", err)
		return nil, err
	}

	return &resp, err
}

func SendMsgConfirmation(followUpMsg *tgbotapi.Message, bot *tgbotapi.BotAPI) (*tgbotapi.Message, error) {
	msg_confirmed := tgbotapi.NewEditMessageText(followUpMsg.Chat.ID, followUpMsg.MessageID, "Проблема решена, спасибо!")
	new_msg, err := bot.Send(msg_confirmed)
	if err != nil {
		log.Printf("Error sending confirmation message: %s", err)
		return nil, err
	}
	return &new_msg, err
}

func AnswerCallback(update *tgbotapi.Update, bot *tgbotapi.BotAPI, followUpMsg *tgbotapi.Message) error {
	if update.CallbackQuery == nil {
		err := errors.New("update.CallbackQuery is nil")
		return err
	}
	switch update.CallbackQuery.Data {
	case "request_accepted":
		_, err := GetCallbackQueryResponse(update, bot) //send callback query to tgapi and get its response
		if err != nil {
			log.Printf("Error calling GetCallbackQueryResponse: %s", err)
		}
		_, err = EditFollowUpMessage(followUpMsg, bot)

		if err != nil {
			log.Printf("Error calling EditFollowUpMessage: %s", err)
		}
	case "request_satisfied":
		_, err := GetCallbackQueryResponse(update, bot)
		if err != nil {
			log.Printf("Error calling GetCallbackQueryResponse: %s", err)
		}
		new_msg, err := SendMsgConfirmation(followUpMsg, bot)
		if err != nil {
			log.Printf("Error sending confirmation: %s", err)
		}
		log.Printf("telegram returned message: %v", new_msg)
	default:
		log.Printf("No valid callback data found")
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
