package commands

import (
	"fmt"

	"gopkg.in/telegram-bot-api.v4"
)

type BatmanHandler struct {
}

var batmanHandlerInfo = CommandInfo{
	Command:     "batman",
	Args:        "",
	Permission:  3,
	Description: "says who is batman",
	LongDesc:    "",
	Usage:       "/batman",
	Examples: []string{
		"/batman",
	},
	ResType: "message",
}

func (responder BatmanHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprint(bot.Self.FirstName, bot.Self.LastName, "is Batman"))
	bot.Send(msg)
}

func (responder BatmanHandler) Info() *CommandInfo {
	return &batmanHandlerInfo
}
