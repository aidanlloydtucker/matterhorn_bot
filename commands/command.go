package commands

import (
	"regexp"

	"gopkg.in/telegram-bot-api.v4"
)

type CommandInfo struct {
	Command     string
	Args        string
	ArgsRegex   regexp.Regexp
	Permission  int
	Description string
	LongDesc    string
	Usage       string
	Examples    []string
	ResType     string
}

type Command interface {
	Info() *CommandInfo
	HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string)
}
