package commands

import (
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

type LinesHandler struct {
}

var linesHandlerInfo = CommandInfo{
	Command:     "lines",
	Args:        `(.+)`,
	Permission:  3,
	Description: "make word lines",
	LongDesc:    "",
	Usage:       "/lines [word]",
	Examples: []string{
		"/lines hello",
	},
	ResType: "message",
}

func (h LinesHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
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

func (h LinesHandler) Info() *CommandInfo {
	return &linesHandlerInfo
}

func (h LinesHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
