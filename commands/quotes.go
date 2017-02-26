package commands

import (
	"fmt"
	chatpkg "github.com/billybobjoeaglt/matterhorn_bot/chat"
	"gopkg.in/telegram-bot-api.v4"
)

type QuotesHandler struct {
	ds *chatpkg.Datastore
}

var quotesHandlerInfo = CommandInfo{
	Command:     "quotes",
	Args:        "",
	Permission:  3,
	Description: "gets a random quote from the chat quotes document",
	LongDesc:    "",
	Usage:       "/quotes",
	Examples: []string{
		"/quotes",
	},
	ResType: "message",
	Hidden:  false,
}

func (h *QuotesHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig

	quote, err := getQuote(h.ds, message.Chat.ID)
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		if quote.Manual {
			msg = tgbotapi.NewMessage(message.Chat.ID, quote.Text)
		} else {
			msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(`"%s" - %s %s`, quote.Text, quote.Author, quote.Date.Format("01/2/06")))
		}
	}
	bot.Send(msg)
}

func (h *QuotesHandler) Info() *CommandInfo {
	return &quotesHandlerInfo
}

func (h *QuotesHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

/*
Params:
*chat.Datastore datastore (required) // the datastore struct to get quotes
*/
func (h *QuotesHandler) Setup(setupFields map[string]interface{}) {
	if dsVal, ok := setupFields["datastore"]; ok {
		if ds, ok := dsVal.(*chatpkg.Datastore); ok {
			h.ds = ds
		}
	}
}
