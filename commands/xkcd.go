package commands

import (
	"errors"
	"math/rand"

	"net/http"

	"strconv"

	"encoding/json"

	"log"

	"gopkg.in/telegram-bot-api.v4"
)

type XkcdHandler struct {
}

var xkcdHandlerInfo = CommandInfo{
	Command:     "xkcd",
	Args:        `(.*)`,
	Permission:  3,
	Description: "gets xkcd",
	LongDesc:    "",
	Usage:       "/xkcd ('new', 'latest', or the id)",
	Examples: []string{
		"/xkcd",
		"/xkcd new",
		"/xkcd 314",
	},
	ResType: "message",
}

func init() {
	// Get latest XKCD
	post, err := GetXKCD(0)
	if err != nil {
		log.Println("Could not get latest XKCD:", err)
	} else {
		xkcdLatestID = post.ID
	}
}

func (h *XkcdHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig

	var xkcdID int

	if len(args) == 0 || args[0] == "" {
		xkcdID = -1
	} else if args[0] == "new" || args[0] == "latest" {
		xkcdID = 0
	} else {
		conv, err := strconv.Atoi(args[0])
		if err != nil {
			xkcdID = 0
		} else {
			xkcdID = conv
		}
	}

	post, err := GetXKCD(xkcdID)
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, "<b>"+post.Title+"</b>\n───\n<i>"+post.Alt+"</i>\n"+post.Img)
		msg.ParseMode = "HTML"
	}

	bot.Send(msg)
}

func (h *XkcdHandler) Info() *CommandInfo {
	return &xkcdHandlerInfo
}

func (h *XkcdHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

func (h *XkcdHandler) Setup(setupFields map[string]interface{}) {

}

type XKCDPost struct {
	Title string `json:"title"`
	Alt   string `json:"alt"`
	Img   string `json:"img"`
	ID    int    `json:"num"`
}

var xkcdLatestID int

func GetXKCD(id int) (*XKCDPost, error) {
	xkcdStr := "http://xkcd.com/"

	if id != 0 {
		if id == -1 {
			if xkcdLatestID == 0 {
				post, err := GetXKCD(0)
				if err != nil {
					return nil, err
				}
				xkcdLatestID = post.ID
			}
			id = rand.Intn(xkcdLatestID-1) + 1
		}
		xkcdStr += strconv.Itoa(id) + "/"
	}

	xkcdStr += "/info.0.json"

	resp, err := http.Get(xkcdStr)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Invalid Status Code: " + resp.Status)
	}

	post := XKCDPost{}
	err = json.NewDecoder(resp.Body).Decode(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil

}
