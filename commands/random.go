package commands

import (
	"math/rand"
	"strconv"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

type RandomHandler struct {
}

func init() {
	rand.Seed(time.Now().Unix())
}

var randomHandlerInfo = CommandInfo{
	Command:     "random",
	Args:        ` ?(.[^ ]*)? ?(.[^ ]*)?`,
	Permission:  3,
	Description: "gets random ",
	LongDesc:    "",
	Usage:       "/random (min) (max)",
	Examples: []string{
		"/random",
		"/random 100",
		"/random 10 100",
	},
	ResType: "message",
}

func (h RandomHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var errMsg tgbotapi.MessageConfig

	defer func(bot *tgbotapi.BotAPI) {
		if errMsg.Text != "" {
			bot.Send(errMsg)
		}
	}(bot)

	min := 0
	max := 100

	if len(args) > 0 && args[0] != "" {
		arg1, err := strconv.Atoi(args[0])
		if err != nil {
			errMsg = NewErrorMessage(message.Chat.ID, err)
			return
		}
		if len(args) > 1 && args[1] != "" {
			arg2, err := strconv.Atoi(args[1])
			if err != nil {
				errMsg = NewErrorMessage(message.Chat.ID, err)
				return
			}
			min = arg1
			max = arg2
		} else {
			max = arg1
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, strconv.Itoa(rand.Intn(max-min)+min))

	bot.Send(msg)
}

func (h RandomHandler) Info() *CommandInfo {
	return &randomHandlerInfo
}

func (h RandomHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
