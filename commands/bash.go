package commands

import (
	"math/rand"

	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/telegram-bot-api.v4"
)

type BashHandler struct {
}

var bashHandlerInfo = CommandInfo{
	Command:     "bash",
	Args:        "",
	Permission:  3,
	Description: "gets a bash.org quote",
	LongDesc:    "",
	Usage:       "/bash",
	Examples: []string{
		"/bash",
	},
	ResType: "message",
}

func (responder BashHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig

	err, quote := GetBash()
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, quote)
	}
	bot.Send(msg)
}

func (responder BashHandler) Info() *CommandInfo {
	return &bashHandlerInfo
}

func GetBash() (error, string) {
	doc, err := goquery.NewDocument("http://bash.org/?random1")
	if err != nil {
		return err, ""
	}

	rand.Seed(time.Now().UTC().UnixNano())

	qtList := doc.Find(".qt")

	return nil, qtList.Eq(rand.Intn(qtList.Length() - 1)).Text()

}
