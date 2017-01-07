package commands

import (
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

type InfoHandler struct {
}

var infoHandlerInfo = CommandInfo{
	Command:     "info",
	Args:        "",
	Permission:  3,
	Description: "shares info about bot",
	LongDesc:    "",
	Usage:       "/info",
	Examples: []string{
		"/info",
	},
	ResType: "message",
}
var BotInfoVersion string
var BotInfoTimestamp *time.Time

func (h InfoHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var biT string
	if BotInfoTimestamp != nil {
		biT = BotInfoTimestamp.String()
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "<b>"+GetUserTitle(&bot.Self)+"</b>\n"+
		"Bot Version: "+BotInfoVersion+"\n"+
		"Build Timestamp: "+biT+"\n"+
		"Github Repo: https://github.com/billybobjoeaglt/matterhorn_bot")
	msg.DisableWebPagePreview = true
	msg.ParseMode = "HTML"

	bot.Send(msg)
}

func (h InfoHandler) Info() *CommandInfo {
	return &infoHandlerInfo
}

func (h InfoHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
