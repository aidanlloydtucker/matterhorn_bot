package commands

import (
	"strconv"

	"gopkg.in/telegram-bot-api.v4"
)

type SettingsHandler struct {
}

var settingsHandlerInfo = CommandInfo{
	Command:     "settings",
	Args:        ``,
	Permission:  3,
	Description: "enters settings link to change settings",
	LongDesc:    "",
	Usage:       "/settings",
	Examples: []string{
		"/settings",
	},
	ResType: "message",
}

func (responder SettingsHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	url := SettingsURL + strconv.FormatInt(message.Chat.ID, 10)

	return

	msg := tgbotapi.NewMessage(message.Chat.ID, "<b>To edit your settings, please go to the link below:</b>\n<a href=\""+url+"\">"+url+"</a>")
	msg.ParseMode = "HTML"

	bot.Send(msg)
}

func (responder SettingsHandler) Info() *CommandInfo {
	return &settingsHandlerInfo
}

var SettingsURL string
