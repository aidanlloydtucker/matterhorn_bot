package commands

import (
	"time"

	"strconv"

	"gopkg.in/telegram-bot-api.v4"
)

type BenchHandler struct {
}

func (responder BenchHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, strconv.FormatInt(time.Now().UnixNano(), 10))
	bot.Send(msg)
	return nil
}

func (responder BenchHandler) Info() *CommandInfo {
	return &CommandInfo{
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
}
