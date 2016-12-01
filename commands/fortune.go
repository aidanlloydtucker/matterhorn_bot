package commands

import (
	"os/exec"

	"gopkg.in/telegram-bot-api.v4"
)

type FortuneHandler struct {
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

func (h FortuneHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig

	err, fortune := GetFortune()
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, fortune)
	}
	bot.Send(msg)
}

func (h FortuneHandler) Info() *CommandInfo {
	return &fortuneHandlerInfo
}

func (h FortuneHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

func GetFortune() (error, string) {
	fOut, err := exec.Command("/usr/games/fortune", "-a", "fortunes", "riddles").Output()
	if err != nil {
		return err, ""
	}
	return nil, string(fOut)
}
