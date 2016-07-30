package commands

import (
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

type ClearHandler struct {
}

var clearHandlerInfo = CommandInfo{
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

func (h ClearHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "-_-"+strings.Repeat("\n", 80)+"-_-")
	bot.Send(msg)
}

func (h ClearHandler) Info() *CommandInfo {
	return &clearHandlerInfo
}

func (h ClearHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
