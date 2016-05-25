package commands

import (
	"errors"

	"net/http"

	"gopkg.in/telegram-bot-api.v4"
)

type CatHandler struct {
}

var CatHandlerInfo = CommandInfo{
	Command:     "cat",
	Args:        "",
	Permission:  3,
	Description: "gets a cat photo",
	LongDesc:    "",
	Usage:       "/cat",
	Examples: []string{
		"/cat",
	},
	ResType: "message",
}

func (responder CatHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) error {
	err, photo := GetCat()
	if err != nil {
		msg := NewErrorMessage(message.Chat.ID, err)
		bot.Send(msg)
		return nil
	}
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, photo)

	bot.Send(msg)
	return nil
}

func (responder CatHandler) Info() *CommandInfo {
	return &CatHandlerInfo
}

func GetCat() (error, tgbotapi.FileReader) {
	resp, err := http.Get("http://thecatapi.com/api/images/get?type=jpg")
	if err != nil {
		return err, tgbotapi.FileReader{}
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Invalid Status Code: " + resp.Status), tgbotapi.FileReader{}
	}

	return nil, tgbotapi.FileReader{
		Name:   "cat.jpg",
		Reader: resp.Body,
		Size:   -1,
	}

}
