package commands

import (
	chatpkg "github.com/billybobjoeaglt/matterhorn_bot/chat"
	"gopkg.in/telegram-bot-api.v4"
	"strconv"
)

type QuoteslinkHandler struct {
	url string
	ds  *chatpkg.Datastore
}

var quoteslinkHandlerInfo = CommandInfo{
	Command:     "quoteslink",
	Args:        "",
	Permission:  3,
	Description: "gets link to quotes document",
	LongDesc:    "",
	Usage:       "/quoteslink",
	Examples: []string{
		"/quoteslink",
	},
	ResType: "message",
	Hidden:  false,
}

func (h *QuoteslinkHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	chat, _, err := h.ds.GetChat(message.Chat.ID)
	if err != nil {
		msg := NewErrorMessage(message.Chat.ID, err)
		bot.Send(msg)
		return
	}

	url := h.url + strconv.Itoa(chat.Settings.QuotesDoc)

	msg := tgbotapi.NewMessage(message.Chat.ID, "<b>To view the quotes doc, please go to the link below:</b>\n<a href=\""+url+"\">"+url+"</a>")
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func (h *QuoteslinkHandler) Info() *CommandInfo {
	return &quoteslinkHandlerInfo
}

func (h *QuoteslinkHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

/*
Params:
string url (default: "unknown") // The base URL to access the quotes link webpage
*chat.Datastore datastore (required) // the datastore struct to get quotes
*/
func (h *QuoteslinkHandler) Setup(setupFields map[string]interface{}) {
	h.url = "unknown"
	if urlVal, ok := setupFields["url"]; ok {
		if url, ok := urlVal.(string); ok {
			h.url = url
		}
	}

	if dsVal, ok := setupFields["datastore"]; ok {
		if ds, ok := dsVal.(*chatpkg.Datastore); ok {
			h.ds = ds
		}
	}
}
