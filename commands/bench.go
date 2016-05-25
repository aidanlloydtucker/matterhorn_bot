package commands

import (
	"time"

	"strconv"

	"gopkg.in/telegram-bot-api.v4"
)

type BenchHandler struct {
}

var benchHandlerInfo = CommandInfo{
	Command:     "bench",
	Args:        "",
	Permission:  3,
	Description: "gets unix nano timestamp",
	LongDesc:    "",
	Usage:       "/bench",
	Examples: []string{
		"/bench",
	},
	ResType: "message",
}

func (responder BenchHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, strconv.FormatInt(time.Now().UnixNano(), 10))
	bot.Send(msg)
}

func (responder BenchHandler) Info() *CommandInfo {
	return &benchHandlerInfo
}
