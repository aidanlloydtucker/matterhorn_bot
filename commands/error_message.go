package commands

import "gopkg.in/telegram-bot-api.v4"

func NewErrorMessage(chatID int64, err error) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Error! "+err.Error())
}
