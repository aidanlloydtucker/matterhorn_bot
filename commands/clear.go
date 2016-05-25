package commands

import (
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

type ClearHandler struct {
}

var ClearHandlerInfo = CommandInfo{
	Command:     "clear",
	Args:        "",
	Permission:  3,
	Description: "clears screen",
	LongDesc:    "",
	Usage:       "/clear",
	Examples: []string{
		"/clear",
	},
	ResType: "message",
}

func (responder ClearHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "-_-"+strings.Repeat("\n", 80)+"-_-")
	bot.Send(msg)
}

func (responder ClearHandler) Info() *CommandInfo {
	return &ClearHandlerInfo
}
