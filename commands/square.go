package commands

import (
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

type SquareHandler struct {
}

var squareHandlerInfo = CommandInfo{
	Command:     "square",
	Args:        `(.+)`,
	Permission:  3,
	Description: "squares a word",
	LongDesc:    "",
	Usage:       "/square [word]",
	Examples: []string{
		"/square hello",
	},
	ResType: "message",
}

func (responder SquareHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	word := args[0]
	word = strings.ToUpper(word)
	var sendStr string
	for _, char := range word {
		sendStr += string(char) + " "
	}
	for _, char := range word[1:] {
		sendStr += "\n" + string(char)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, sendStr)
	bot.Send(msg)
}

func (responder SquareHandler) Info() *CommandInfo {
	return &squareHandlerInfo
}
