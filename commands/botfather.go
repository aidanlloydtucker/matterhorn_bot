package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

type BotFatherHandler struct {
}

var botFatherHandlerInfo = CommandInfo{
	Command:     "botfather",
	Args:        "",
	Permission:  3,
	Description: "gets botfather list",
	LongDesc:    "",
	Usage:       "/botfather",
	Examples: []string{
		"/botfather",
	},
	ResType: "message",
	Hidden:  true,
}

func (h BotFatherHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msgStr string
	for _, cmd := range CommandMap {
		if cmd.Hidden {
			continue
		}
		msgStr += cmd.Command + " - " + cmd.Description + "\n"
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgStr)
	bot.Send(msg)
}

func (h BotFatherHandler) Info() *CommandInfo {
	return &botFatherHandlerInfo
}

func (h BotFatherHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
