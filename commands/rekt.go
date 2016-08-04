package commands

import (
	"io/ioutil"
	"math/rand"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

var reks []string

func init() {
	fileBytes, err := ioutil.ReadFile("./resources/reks.txt")
	if err != nil {
		return
	}
	fileStr := string(fileBytes)
	reks = strings.Split(fileStr, "\n")
}

type RektHandler struct {
}

var rektHandlerInfo = CommandInfo{
	Command:     "rekt",
	Args:        `(.+)`,
	Permission:  3,
	Description: "reks person",
	LongDesc:    "",
	Usage:       "/rekt [name]",
	Examples: []string{
		"/rekt bob",
	},
	ResType: "message",
}

func (h RektHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	res := strings.Replace(reks[rand.Intn(len(reks))], "$USER", args[0], -1)

	msg := tgbotapi.NewMessage(message.Chat.ID, res)
	bot.Send(msg)
}

func (h RektHandler) Info() *CommandInfo {
	return &rektHandlerInfo
}

func (h RektHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return true, message.ReplyToMessage.From.String()
}
