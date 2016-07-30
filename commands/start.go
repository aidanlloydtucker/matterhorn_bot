package commands

import "gopkg.in/telegram-bot-api.v4"

type StartHandler struct {
}

var startHandlerInfo = CommandInfo{
	Command:     "start",
	Args:        ``,
	Permission:  3,
	Description: "start message",
	LongDesc:    "",
	Usage:       "/start",
	Examples: []string{
		"/start",
	},
	ResType: "message",
}

func (h StartHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {

	msg := tgbotapi.NewMessage(message.Chat.ID, "<b>Hello and welcome to "+bot.Self.UserName+"!</b>\n---\nTo setup your chat, type /settings\nTo look at "+bot.Self.UserName+"'s many commands, type /help")
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func (h StartHandler) Info() *CommandInfo {
	return &startHandlerInfo
}

func (h StartHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
