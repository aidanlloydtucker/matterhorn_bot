package commands

import (
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

type InfoHandler struct {
	botVersion      string
	botTimestamp    *time.Time
	botTimestampStr string
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

func (h *InfoHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "<b>"+GetUserTitle(&bot.Self)+"</b>\n"+"Bot Version: "+h.botVersion+"\n"+"Build Timestamp: "+h.botTimestampStr+"\n"+
		"Github Repo: https://github.com/billybobjoeaglt/matterhorn_bot")
	msg.DisableWebPagePreview = true
	msg.ParseMode = "HTML"

	bot.Send(msg)
}

func (h *InfoHandler) Info() *CommandInfo {
	return &infoHandlerInfo
}

func (h *InfoHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

/*
Params:
string botVersion (default: "unknown") // Version of bot
*time.Time botTimestamp (default: nil) // Time the bot was built
*/
func (h *InfoHandler) Setup(setupFields map[string]interface{}) {
	h.botVersion = "unknown"
	h.botTimestamp = nil

	if versionVal, ok := setupFields["botVersion"]; ok {
		if version, ok := versionVal.(string); ok {
			if version != "" {
				h.botVersion = version
			}
		}
	}

	if timestampVal, ok := setupFields["botTimestamp"]; ok {
		if timestamp, ok := timestampVal.(*time.Time); ok {
			h.botTimestamp = timestamp
		}
	}

	if h.botTimestamp != nil {
		h.botTimestampStr = h.botTimestamp.String()
	} else {
		h.botTimestampStr = "unknown"
	}
}
