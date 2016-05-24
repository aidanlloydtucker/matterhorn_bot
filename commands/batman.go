package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

type BatmanHandler struct {

}

func (responder BatmanHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Sansa Stark is Batman")
	bot.Send(msg)
	return nil
}

func (responder BatmanHandler) Info() *CommandInfo {
	return &CommandInfo{
		Command: "batman",
		Args: "",
		Permission: 3,
		Description: "says who is batman",
		LongDesc: "",
		Usage: "/batman",
		Examples: []string{
			"/batman",
		},
		ResType: "message",
	}
}