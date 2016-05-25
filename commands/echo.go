package commands

import "gopkg.in/telegram-bot-api.v4"

type EchoHandler struct {
}

var EchoHandlerInfo = CommandInfo{
	Command:     "echo",
	Args:        `(.+)`,
	Permission:  3,
	Description: "echos input",
	LongDesc:    "",
	Usage:       "/echo [input]",
	Examples: []string{
		"/echo hello world",
	},
	ResType: "message",
}

func (responder EchoHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, args[0])
	bot.Send(msg)
}

func (responder EchoHandler) Info() *CommandInfo {
	return &EchoHandlerInfo
}
