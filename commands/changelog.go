package commands

import (
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"strings"
)

type ChangelogHandler struct {
	changes string
}

var changelogHandlerInfo = CommandInfo{
	Command:     "changelog",
	Args:        "",
	Permission:  3,
	Description: "",
	LongDesc:    "",
	Usage:       "/changelog",
	Examples: []string{
		"/changelog",
	},
	ResType: "message",
	Hidden:  false,
}

func (h *ChangelogHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "*Changelog*\n---\n"+h.changes)
	msg.ParseMode = "markdown"
	bot.Send(msg)
}

func (h *ChangelogHandler) Info() *CommandInfo {
	return &changelogHandlerInfo
}

func (h *ChangelogHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

/*
Params:
string path (default: ./changelog.txt) // path to the file of changes
*/
func (h *ChangelogHandler) Setup(setupFields map[string]interface{}) {
	path := "./changelog.txt"
	if pathVal, ok := setupFields["path"]; ok {
		if newPath, ok := pathVal.(string); ok {
			path = newPath
		}
	}
	if path != "" {
		fileBytes, err := ioutil.ReadFile(path)
		if err == nil {
			fileStr := string(fileBytes)
			h.changes = strings.Join(strings.Split(fileStr, "\n\n")[:5], "\n\n")
		}
	}
}
