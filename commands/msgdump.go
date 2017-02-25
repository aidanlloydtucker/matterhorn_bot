package commands

import (
	"encoding/json"
	"gopkg.in/telegram-bot-api.v4"
)

type MsgDumpHandler struct {
}

var msgDumpHandlerInfo = CommandInfo{
	Command:     "msgdump",
	Args:        "",
	Permission:  3,
	Description: "dumps the info recieved in the message",
	LongDesc:    "",
	Usage:       "/msgdump",
	Examples: []string{
		"/msgdump",
	},
	ResType: "message",
}

func (h *MsgDumpHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig

	jsonBytes, err := json.MarshalIndent(message, "", "	")
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, string(jsonBytes))
	}

	bot.Send(msg)
}

func (h *MsgDumpHandler) Info() *CommandInfo {
	return &msgDumpHandlerInfo
}

func (h *MsgDumpHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

func (h *MsgDumpHandler) Setup(setupFields map[string]interface{}) {

}
