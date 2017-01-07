package commands

import (
	"math/rand"
	"regexp"
	"time"

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
	Hidden      bool
}

type Command interface {
	Info() *CommandInfo
	HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string)
	HandleReply(message *tgbotapi.Message) (bool, string)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewErrorMessage(chatID int64, err error) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Error! "+err.Error())
}

func GetUserTitle(user *tgbotapi.User) string {
	name := user.FirstName
	if user.LastName != "" {
		name += " " + user.LastName
	}

	return name
}
