package test

import (
	"ingest_bot/helpers"
	"log"
	te "testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestCallIngestCmd(t *te.T) {
	var request_msg tgbotapi.Message
	request_msg.Chat.ID = 4710669718
	request_msg.Text = "привет инжест!"
	request_msg.MessageID = 12
	fwd_msg := tgbotapi.NewForward(helpers.Settings.IngestChatId, request_msg.Chat.ID, request_msg.MessageID)
	log.Println(fwd_msg)
	/*_, err := bot.Send(fwd_msg)
	if err != nil {
		log.Printf("Request message could not be forwarded, error: %s", err)

	}*/
}
