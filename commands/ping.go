package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

type PingHandler struct {
}

var pingHandlerInfo = CommandInfo{
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

func (h PingHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "PONG!")
	bot.Send(msg)
}

func (h PingHandler) Info() *CommandInfo {
	return &pingHandlerInfo
}

func (h PingHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
