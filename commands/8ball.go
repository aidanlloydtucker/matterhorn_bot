package commands

import (
	"math/rand"

	"gopkg.in/telegram-bot-api.v4"
)

var magicBallAnswers = []string{
	"It is certain",
	"It is decidedly so",
	"Without a doubt",
	"Yes definitely",
	"You may rely on it",
	"As I see it yes",
	"Most likely",
	"Outlook good",
	"Yes",
	"Signs point to yes",
	"Reply hazy try again",
	"Ask again later",
	"Better not tell you now",
	"Cannot predict now",
	"Concentrate and ask again",
	"Don't count on it",
	"My reply is no",
	"My sources say no",
	"Outlook not so good",
	"Very doubtful",
}

type MagicBallHandler struct {
}

var magicBallHandlerInfo = CommandInfo{
	Command:     "8ball",
	Permission:  3,
	Description: "gets a response from the Magic 8-Ball™",
	LongDesc:    "",
	Usage:       "/8ball",
	Examples: []string{
		"/8ball",
	},
	ResType: "message",
}

func (h MagicBallHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, magicBallAnswers[rand.Intn(len(magicBallAnswers))])
	bot.Send(msg)
}

func (h MagicBallHandler) Info() *CommandInfo {
	return &magicBallHandlerInfo
}

func (h MagicBallHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}
