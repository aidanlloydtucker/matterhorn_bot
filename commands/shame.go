package commands

import "gopkg.in/telegram-bot-api.v4"

type ShameHandler struct {
}

var shameHandlerInfo = CommandInfo{
	Command:     "shame",
	Args:        "",
	Permission:  3,
	Description: "SHAME! SHAME! SHAME! DING DING!",
	LongDesc:    "",
	Usage:       "/shame",
	Examples: []string{
		"/shame",
	},
	ResType: "message",
}

var shameFileID string

func (h ShameHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.VoiceConfig

	if shameFileID != "" {
		msg = tgbotapi.NewVoiceShare(message.Chat.ID, shameFileID)
	} else {
		msg = tgbotapi.NewVoiceUpload(message.Chat.ID, "./resources/shame_sfx.mp3")
	}
	retMsg, _ := bot.Send(msg)
	if retMsg.Voice.FileID != "" {
		shameFileID = retMsg.Voice.FileID
	}

}

func (h ShameHandler) Info() *CommandInfo {
	return &shameHandlerInfo
}

func (h ShameHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
