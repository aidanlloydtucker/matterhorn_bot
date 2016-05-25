package commands

import "gopkg.in/telegram-bot-api.v4"

type HelpHandler struct {
}

var helpHandlerInfo = CommandInfo{
	Command:     "help",
	Args:        "",
	Permission:  3,
	Description: "lists commands",
	LongDesc:    "",
	Usage:       "/help",
	Examples: []string{
		"/help",
	},
	ResType: "message",
}

func (responder HelpHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msgStr string = "<b>Commands:</b>\n"
	for _, cmd := range *CommandList {
		msgStr += "• " + cmd.Info().Command + " - " + cmd.Info().Description + "\n"
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgStr)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func (responder HelpHandler) Info() *CommandInfo {
	return &helpHandlerInfo
}

var CommandList *[]Command
