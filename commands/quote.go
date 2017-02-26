package commands

import (
	chatpkg "github.com/billybobjoeaglt/matterhorn_bot/chat"
	"gopkg.in/telegram-bot-api.v4"
	"strings"
	"time"
)

type QuoteHandler struct {
	ds *chatpkg.Datastore
}

var quoteHandlerInfo = CommandInfo{
	Command:     "quote",
	Args:        "",
	Permission:  3,
	Description: "adds a quote to the chat quotes document",
	LongDesc:    "it can be used as a reply to a message to automatically quote it (even works with forwarded messages!)",
	Usage:       "/quote [quote]",
	Examples: []string{
		"/quote \"Hello World\" - MatterhornBot 2017",
	},
	ResType: "message",
	Hidden:  false,
}

const quoteReplyInfoPrefix = "/REPLY/ "
const quoteSuccessMessage = "Added Quote!"

func (h *QuoteHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	if strings.HasPrefix(message.CommandArguments(), quoteReplyInfoPrefix) {
		msg := tgbotapi.NewMessage(message.Chat.ID, strings.TrimPrefix(message.CommandArguments(), quoteReplyInfoPrefix))
		bot.Send(msg)
		return
	}

	quote := Quote{
		Text:   message.CommandArguments(),
		Manual: true,
	}

	err := addQuote(h.ds, message.Chat.ID, quote)
	if err != nil {
		msg := NewErrorMessage(message.Chat.ID, err)
		bot.Send(msg)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, quoteSuccessMessage)
	bot.Send(msg)
}

func (h *QuoteHandler) Info() *CommandInfo {
	return &quoteHandlerInfo
}

// TODO: Figure out a better way to send a message than doing "/REPLY/"
func (h *QuoteHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	if message.ReplyToMessage.Text == "" {
		return false, ""
	}

	var fromUser *tgbotapi.User
	var fromDate int

	if message.ReplyToMessage.ForwardFrom != nil {
		fromUser = message.ReplyToMessage.ForwardFrom
		fromDate = message.ReplyToMessage.ForwardDate
	} else {
		fromUser = message.ReplyToMessage.From
		fromDate = message.ReplyToMessage.Date
	}

	var name string

	if fromUser.FirstName != "" {
		name = fromUser.FirstName
		if message.ReplyToMessage.From.LastName != "" {
			name += " " + fromUser.LastName
		}
	} else if fromUser.UserName != "" {
		name = fromUser.UserName
	} else {
		name = "Unknown"
	}

	var timeSent time.Time
	if fromDate != 0 {
		timeSent = time.Unix(int64(fromDate), 0)
	} else {
		timeSent = time.Now()
	}

	quote := Quote{
		Text:   message.ReplyToMessage.Text,
		Date:   timeSent,
		Author: name,
		Manual: false,
	}

	err := addQuote(h.ds, message.Chat.ID, quote)
	if err != nil {
		return true, quoteReplyInfoPrefix + "Error! " + err.Error()
	}

	return true, quoteReplyInfoPrefix + quoteSuccessMessage
}

/*
Params:
*chat.Datastore datastore (required) // the datastore struct to get quotes
*/
func (h *QuoteHandler) Setup(setupFields map[string]interface{}) {
	if dsVal, ok := setupFields["datastore"]; ok {
		if ds, ok := dsVal.(*chatpkg.Datastore); ok {
			h.ds = ds
		}
	}
}
