package commands

import (
	"gopkg.in/telegram-bot-api.v4"
	"fmt"
)

type ChatidHandler struct {
}

var chatidHandlerInfo = CommandInfo{
	Command:     "chatid",
	Args:        "",
	Permission:  3,
	Description: "gets the chat's chat ID",
	LongDesc:    "",
	Usage:       "/chatid",
	Examples: []string{
		"/chatid",
	},
	ResType: "message",
}

func (h ChatidHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("This chat's chat ID is: %d", message.Chat.ID))
	bot.Send(msg)
}

func (h ChatidHandler) Info() *CommandInfo {
	return &chatidHandlerInfo
}

func (h ChatidHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
