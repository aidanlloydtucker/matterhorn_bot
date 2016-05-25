package commands

import (
	"errors"

	"net/http"

	"encoding/json"

	"strconv"

	"gopkg.in/telegram-bot-api.v4"
)

type BitcoinHandler struct {
}

var BitcoinHandlerInfo = CommandInfo{
	Command:     "bitcoin",
	Args:        "",
	Permission:  3,
	Description: "gets unix nano timestamp",
	LongDesc:    "",
	Usage:       "/bitcoin",
	Examples: []string{
		"/bitcoin",
	},
	ResType: "message",
}

func (responder BitcoinHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) error {
	var msg tgbotapi.MessageConfig

	err, price := GetBitcoin()
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Bitcoin most recent price: $"+price)
	}
	bot.Send(msg)
	return nil
}

func (responder BitcoinHandler) Info() *CommandInfo {
	return &BitcoinHandlerInfo
}

func GetBitcoin() (error, string) {
	resp, err := http.Get("http://api.bitcoinaverage.com/all")
	if err != nil {
		return err, ""
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Invalid Status Code: " + resp.Status), ""
	}
	var jsonRes map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonRes)
	if err != nil {
		return err, ""
	}
	price := strconv.FormatFloat(jsonRes["USD"].(map[string]interface{})["averages"].(map[string]interface{})["last"].(float64), 'f', 2, 64)
	if price == "" {
		return errors.New("Missing Price"), ""
	}
	return nil, price

}
