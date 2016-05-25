package commands

import (
	"net/url"

	"gopkg.in/telegram-bot-api.v4"
)

type LmgtfyHandler struct {
}

var LmgtfyHandlerInfo = CommandInfo{
	Command:     "lmgtfy",
	Args:        "(.+)",
	Permission:  3,
	Description: "let me google that for you",
	LongDesc:    "",
	Usage:       "/lmgtfy [input]",
	Examples: []string{
		"/lmgtfy hello world",
	},
	ResType: "message",
}

func (responder LmgtfyHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "http://lmgtfy.com/?q="+url.QueryEscape(args[0]))
	bot.Send(msg)
}

func (responder LmgtfyHandler) Info() *CommandInfo {
	return &LmgtfyHandlerInfo
}
