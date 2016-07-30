package commands

import "gopkg.in/telegram-bot-api.v4"

type SquareHandler struct {
}

var squareHandlerInfo = CommandInfo{
	Command:     "square",
	Args:        `(.+)`,
	Permission:  3,
	Description: "square a word",
	LongDesc:    "",
	Usage:       "/square [word]",
	Examples: []string{
		"/square hello",
	},
	ResType: "message",
}

func (h SquareHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	word := args[0]
	var sendStr string

	for i := 0; i < len(word); i++ {
		sendStr += word[i:] + word[:i] + "\n"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, sendStr)
	bot.Send(msg)
}

func (h SquareHandler) Info() *CommandInfo {
	return &squareHandlerInfo
}

func (h SquareHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
