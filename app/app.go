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
		helpers.SendMeInfo(err.Error(), bot)
		return err
	}
	return err
}

func EditFollowUpMessage(respChatID int64, respMessageID int, bot *tgbotapi.BotAPI) (*tgbotapi.Message, error) {
	// if resp == nil {
	// 	log.Panic("followUpMsg is nil")

	// }
	// if resp.Chat == nil {
	// 	log.Panic("followUpMsg.Chat is nil")
	// }
	editedKeyboard := tgbotapi.NewEditMessageReplyMarkup(respChatID, respMessageID, ConfirmKeyboard)
	newResp, err := bot.Send(editedKeyboard)
	if err != nil {
		log.Printf("Error while sending editedKeyboard to telegram : %s", err)
		return nil, err
	}

	return &newResp, err
}

func SendMsgConfirmation(respChatID int64, respMessageID int, bot *tgbotapi.BotAPI) (*tgbotapi.Message, error) {
	msg_confirmed := tgbotapi.NewEditMessageText(respChatID, respMessageID, "Проблема решена, спасибо!")
	new_msg, err := bot.Send(msg_confirmed)
	if err != nil {
		log.Printf("Error sending confirmation message: %s", err)
		return nil, err
	}
	return &new_msg, err
}

func AnswerCallback(newUpd *tgbotapi.Update, bot *tgbotapi.BotAPI, followUpMsg *tgbotapi.Message) error {
	if newUpd.CallbackQuery == nil {
		err := errors.New("update.CallbackQuery is nil")
		return err
	}
	switch newUpd.CallbackQuery.Data {
	case "request_accepted":
		err := GetCallbackQueryResponse(newUpd, bot) //send callback query to tgapi
		//callback := update                              //передать месседж айди для сообщения которое нужно поменять через колбек
		if err != nil {
			log.Printf("Error calling GetCallbackQueryResponse: %s", err)
		}
		if newUpd == nil {
			log.Panic("newUpd is nil")
		}
		if newUpd.CallbackQuery == nil {
			log.Panic("newUpd.Message is nil")
		}
		_, err = EditFollowUpMessage(newUpd.CallbackQuery.Message.Chat.ID, newUpd.CallbackQuery.Message.MessageID, bot)

		if err != nil {
			log.Printf("Error calling EditFollowUpMessage: %s", err)
		}
	case "request_satisfied":
		err := GetCallbackQueryResponse(newUpd, bot)
		if err != nil {
			log.Printf("Error calling GetCallbackQueryResponse: %s", err)
		}
		if newUpd == nil {
			log.Panic("newUpd is nil")
		}
		if newUpd.CallbackQuery == nil {
			log.Panic("newUpd.Message is nil")
		}
		new_msg, err := SendMsgConfirmation(newUpd.CallbackQuery.Message.Chat.ID, newUpd.CallbackQuery.Message.MessageID, bot)
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
