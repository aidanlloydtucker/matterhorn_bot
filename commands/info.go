package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

type InfoHandler struct {
}

var infoHandlerInfo = CommandInfo{
	Command:     "info",
	Args:        "",
	Permission:  3,
	Description: "shares info about bot",
	LongDesc:    "",
	Usage:       "/info",
	Examples: []string{
		"/info",
	},
	ResType: "message",
}

func (h InfoHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "This bot was created by Aidan Lloyd-Tucker (telegram: @slaidan_lt)\n"+
		"The github repo is: https://github.com/billybobjoeaglt/matterhorn_bot")
	bot.Send(msg)
}

func (h InfoHandler) Info() *CommandInfo {
	return &infoHandlerInfo
}

func (h InfoHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
