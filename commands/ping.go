package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

type PingHandler struct {
}

var PingHandlerInfo = CommandInfo{
	Command:     "ping",
	Args:        "",
	Permission:  3,
	Description: "pings bot",
	LongDesc:    "",
	Usage:       "/ping",
	Examples: []string{
		"/ping",
	},
	ResType: "message",
}

func (responder PingHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "PONG!")
	bot.Send(msg)
}

func (responder PingHandler) Info() *CommandInfo {
	return &PingHandlerInfo
}
