package main

import (
	"github.com/billybobjoeaglt/sansa_bot/commands"
	"gopkg.in/telegram-bot-api.v4"
)

type Command interface {
	Info() *commands.CommandInfo
	HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error
}
