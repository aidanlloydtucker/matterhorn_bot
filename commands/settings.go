package commands

import (
	"strconv"

	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type SettingsHandler struct {
	url string
}

var settingsHandlerInfo = CommandInfo{
	Command:     "settings",
	Args:        ``,
	Permission:  3,
	Description: "enters settings link to change settings",
	LongDesc:    "",
	Usage:       "/settings",
	Examples: []string{
		"/settings",
	},
	ResType: "message",
}

func (h *SettingsHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	url := h.url + strconv.FormatInt(message.Chat.ID, 10)

	msg := tgbotapi.NewMessage(message.Chat.ID, "<b>To edit your settings, please go to the link below:</b>\n<a href=\""+url+"\">"+url+"</a>")
	msg.ParseMode = "HTML"

	bot.Send(msg)
}

func (h *SettingsHandler) Info() *CommandInfo {
	return &settingsHandlerInfo
}

func (h *SettingsHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

/*
Params:
string url (default: "unknown") // The base URL to access the settings webpage
*/
func (h *SettingsHandler) Setup(setupFields map[string]interface{}) {
	log.Println("recieved")
	h.url = "unknown"

	if urlVal, ok := setupFields["url"]; ok {
		log.Println("ok", urlVal)
		if url, ok := urlVal.(string); ok {
			log.Println("okok", url)
			h.url = url
		}
	}
	log.Println("h:", h)
}
