package commands

import "gopkg.in/telegram-bot-api.v4"

type EchoHandler struct {
}

var echoHandlerInfo = CommandInfo{
	Command:     "echo",
	Args:        `(.+)`,
	Permission:  3,
	Description: "echos input",
	LongDesc:    "",
	Usage:       "/echo [input]",
	Examples: []string{
		"/echo hello world",
	},
	ResType: "message",
}

func (h *EchoHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, args[0])
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func (h *EchoHandler) Info() *CommandInfo {
	return &echoHandlerInfo
}

func (h *EchoHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

func (h *EchoHandler) Setup(setupFields map[string]interface{}) {

}
