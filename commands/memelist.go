package commands

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"errors"

	"gopkg.in/telegram-bot-api.v4"
)

var memeList string

func getMemeList() string {
	files, err := ioutil.ReadDir("./resources/meme-tmpl/")
	if err != nil {
		return ""
	}

	var msgStr string

	for _, f := range files {
		msgStr += "• " + strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())) + "\n"
	}
	return msgStr
}

func init() {
	memeList = getMemeList()
}

type MemeListHandler struct {
}

var memeListHandlerInfo = CommandInfo{
	Command:     "memelist",
	Args:        ``,
	Permission:  3,
	Description: "lists all available memes",
	LongDesc:    "",
	Usage:       `/memelist`,
	Examples: []string{
		`/memelist`,
	},
	ResType: "message",
}

func (h *MemeListHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var errMsg tgbotapi.MessageConfig

	defer func(bot *tgbotapi.BotAPI) {
		if errMsg.Text != "" {
			bot.Send(errMsg)
		}
	}(bot)

	var msgMemeLs string
	if memeList == "" {
		msgMemeLs = getMemeList()
	} else {
		msgMemeLs = memeList
	}

	if msgMemeLs == "" {
		errMsg = NewErrorMessage(message.Chat.ID, errors.New("Couldn't find memes"))
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "<b>Meme List:</b>\n"+msgMemeLs)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func (h *MemeListHandler) Info() *CommandInfo {
	return &memeListHandlerInfo
}

func (h *MemeListHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

//TODO: path param
func (h *MemeListHandler) Setup(setupFields map[string]interface{}) {

}
