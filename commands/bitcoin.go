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

var bitcoinHandlerInfo = CommandInfo{
	Command:     "bitcoin",
	Args:        "",
	Permission:  3,
	Description: "gets bitcoin prices in USD",
	LongDesc:    "",
	Usage:       "/bitcoin",
	Examples: []string{
		"/bitcoin",
	},
	ResType: "message",
}

func (h *BitcoinHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig

	err, price := GetBitcoin()
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Bitcoin's most recent price is $"+price)
	}
	bot.Send(msg)
}

func (h *BitcoinHandler) Info() *CommandInfo {
	return &bitcoinHandlerInfo
}

func (h *BitcoinHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

func (h *BitcoinHandler) Setup(setupFields map[string]interface{}) {

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
	var jsonRes = struct {
		USD struct {
			Averages struct {
				Last float64 `json:"last"`
			} `json:"averages"`
		} `json:"USD"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&jsonRes)
	if err != nil {
		return err, ""
	}
	price := strconv.FormatFloat(jsonRes.USD.Averages.Last, 'f', 2, 64)
	if price == "" {
		return errors.New("Missing Price"), ""
	}
	return nil, price

}
