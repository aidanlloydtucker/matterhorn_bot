package commands

import (
	"os/exec"

	"gopkg.in/telegram-bot-api.v4"
)

type FortuneHandler struct {
	path string
}

var fortuneHandlerInfo = CommandInfo{
	Command:     "fortune",
	Args:        "",
	Permission:  3,
	Description: "reads a unix fortune",
	LongDesc:    "",
	Usage:       "/fortune",
	Examples: []string{
		"/fortune",
	},
	ResType: "message",
}

func (h *FortuneHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig

	err, fortune := GetFortune(h.path)
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, fortune)
	}
	bot.Send(msg)
}

func (h *FortuneHandler) Info() *CommandInfo {
	return &fortuneHandlerInfo
}

func (h *FortuneHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

/*
Params:
string path (default: /usr/games/fortune) // Path to fortune command
*/
func (h *FortuneHandler) Setup(setupFields map[string]interface{}) {
	h.path = "/usr/games/fortune"

	if val, ok := setupFields["path"]; ok {
		if path, ok := val.(string); ok {
			h.path = path
		}
	}
}

func GetFortune(path string) (error, string) {
	fOut, err := exec.Command(path, "-a", "fortunes", "riddles").Output()
	if err != nil {
		return err, ""
	}
	return nil, string(fOut)
}
