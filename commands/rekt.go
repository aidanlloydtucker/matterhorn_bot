package commands

import (
	"io/ioutil"
	"math/rand"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

type RektHandler struct {
	reks []string
}

var rektHandlerInfo = CommandInfo{
	Command:     "rekt",
	Args:        `(.+)`,
	Permission:  3,
	Description: "reks person",
	LongDesc:    "",
	Usage:       "/rekt [name]",
	Examples: []string{
		"/rekt bob",
	},
	ResType: "message",
}

func (h *RektHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	res := strings.Replace(h.reks[rand.Intn(len(h.reks))], "$USER", args[0], -1)

	msg := tgbotapi.NewMessage(message.Chat.ID, res)
	bot.Send(msg)
}

func (h *RektHandler) Info() *CommandInfo {
	return &rektHandlerInfo
}

func (h *RektHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	if message.CommandArguments() != "" {
		return true, message.ReplyToMessage.From.String()
	}
	return false, ""
}

/*
Params:
string path (default: ./resources/reks.txt) // path to a file of reks. if blank (""), it will not look up the file and instead use the `reks` param
[]string reks (optional) // Reks to either be appended to the main list of reks or become the list if the path param is blank
*/
func (h *RektHandler) Setup(setupFields map[string]interface{}) {
	h.reks = []string{}

	path := "./resources/reks.txt"
	if pathVal, ok := setupFields["path"]; ok {
		if newPath, ok := pathVal.(string); ok {
			path = newPath
		}
	}
	if path != "" {
		fileBytes, err := ioutil.ReadFile(path)
		if err == nil {
			fileStr := string(fileBytes)
			h.reks = strings.Split(fileStr, "\n")
		}
	}

	if reksVal, ok := setupFields["reks"]; ok {
		if reks, ok := reksVal.([]string); ok {
			h.reks = append(h.reks, reks...)
		}
	}
}
