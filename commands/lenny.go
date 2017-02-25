package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

type LennyHandler struct {
}

var lennyHandlerInfo = CommandInfo{
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

func (h *LennyHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "( ͡° ͜ʖ ͡°)")
	if message.ReplyToMessage != nil {
		msg.ReplyToMessageID = message.ReplyToMessage.MessageID
	} else {
		msg.ReplyToMessageID = message.MessageID
	}
	bot.Send(msg)
}

func (h *LennyHandler) Info() *CommandInfo {
	return &lennyHandlerInfo
}

func (h *LennyHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

func (h *LennyHandler) Setup(setupFields map[string]interface{}) {

}
