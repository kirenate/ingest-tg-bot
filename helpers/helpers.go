package helpers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	ID   int64  `json:"id"`
	Msg  string `json:"msg"`
	Name string `json:"name"`
}

type settings struct {
	Token      string `json:"token"`
	FromChatId int64  `json:"from_chat_id"`
	ToChatId   int64  `json:"to_chat_id"`
	Users      []User `json:"users"`
}

func (r *settings) userIDToUser() map[int64]User {
	users := make(map[int64]User)
	for _, user := range r.Users {
		users[user.ID] = user
	}

	return users
}

func (r *settings) UserNameToID() map[string]int64 {
	users := make(map[string]int64)
	for _, user := range r.Users {
		users[user.Name] = user.ID
	}
	return users
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
	// err = json.Unmarshal(bytes, &Settings.Users)
	// if err != nil {
	// 	panic(err)
	// }
}

func ChooseCustomMessage(fromUserId int64) string {
	for userId, user := range Settings.userIDToUser() {
		if userId == fromUserId {
			return user.Msg
		}
	}
	return "Проблема решена, спасибо!"
}

func SendMeInfo(info string, bot *tgbotapi.BotAPI) {
	var text string = fmt.Sprintf("В боте что-то сломалось...\n %s", info)
	userNameToID := Settings.UserNameToID()
	msg := tgbotapi.NewMessage(userNameToID["Me"], text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Stack().Err(err).Msg("can't send an error message")
	}
}

func PrintJSON(str string, obj any) {
	updateJSON, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to marshal update")
	}

	fmt.Println(str + string(updateJSON))
}
