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

func (h SettingsHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	url := SettingsURL + strconv.FormatInt(message.Chat.ID, 10)

	msg := tgbotapi.NewMessage(message.Chat.ID, "<b>To edit your settings, please go to the link below:</b>\n<a href=\""+url+"\">"+url+"</a>")
	msg.ParseMode = "HTML"

	bot.Send(msg)
}

func (h SettingsHandler) Info() *CommandInfo {
	return &settingsHandlerInfo
}

func (h SettingsHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

var SettingsURL string
