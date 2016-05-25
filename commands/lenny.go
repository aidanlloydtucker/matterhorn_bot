package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

type LennyHandler struct {
}

var LennyHandlerInfo = CommandInfo{
	Command:     "lenny",
	Args:        "",
	Permission:  3,
	Description: "shows lenny face",
	LongDesc:    "",
	Usage:       "/lenny",
	Examples: []string{
		"/lenny",
	},
	ResType: "message",
}

func (responder LennyHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "( ͡° ͜ʖ ͡°)")
	bot.Send(msg)
}

func (responder LennyHandler) Info() *CommandInfo {
	return &LennyHandlerInfo
}
