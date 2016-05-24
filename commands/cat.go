package commands

import (
	"errors"

	"net/http"

	"log"

	"gopkg.in/telegram-bot-api.v4"
)

type CatHandler struct {
}

func (responder CatHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	log.Println("1")
	err, photo := GetCat()
	log.Println("2")
	if err != nil {
		msg := NewErrorMessage(message.Chat.ID, err)
		bot.Send(msg)
		return nil
	}
	log.Println("3")
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, photo)
	log.Println("4")

	bot.Send(msg)
	log.Println("5")
	return nil
}

func (responder CatHandler) Info() *CommandInfo {
	return &CommandInfo{
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
}

func GetCat() (error, tgbotapi.FileReader) {
	resp, err := http.Get("http://thecatapi.com/api/images/get?type=jpg")
	if err != nil {
		return err, tgbotapi.FileReader{}
	}

	if resp.StatusCode >= 400 {
		return errors.New("Invalid Status Code: " + resp.Status), tgbotapi.FileReader{}
	}

	return nil, tgbotapi.FileReader{
		Name:   "cat.jpg",
		Reader: resp.Body,
		Size:   -1,
	}

}
