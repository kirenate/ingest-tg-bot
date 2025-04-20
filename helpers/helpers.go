package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type settings struct {
	Token        string `json:"token"`
	NozhiChatId  int64  `json:"nozhi_chat_id"`
	IngestChatId int64  `json:"ingest_chat_id"`
	EkirenId     int64  `json:"ekiren_id"`
}

var Settings settings

func init() {
	bytes, err := os.ReadFile("./helpers/.credentials.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &Settings)
	if err != nil {
		panic(err)
	}
}

func SendMeInfo(info string, bot *tgbotapi.BotAPI) {
	var text string = fmt.Sprintf("В боте что-то сломалось...\n %s", info)
	msg := tgbotapi.NewMessage(Settings.EkirenId, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("You can't even get an error message.\n %s", err)
	}
}
