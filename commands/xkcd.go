package commands

import (
	"errors"
	"math/rand"

	"net/http"

	"strconv"

	"encoding/json"

	"time"

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
	Usage:       "/xkcd ('new' or a number)",
	Examples: []string{
		"/xkcd",
		"/xkcd new",
		"/xkcd 314",
	},
	ResType: "message",
}

func (responder XkcdHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig

	var xkcdId int

	if len(args) == 0 || args[0] == "" {
		xkcdId = -1
	} else if args[0] == "new" {
		xkcdId = 0
	} else {
		conv, err := strconv.Atoi(args[0])
		if err != nil {
			xkcdId = 0
		} else {
			xkcdId = conv
		}
	}

	err, xkcd := GetXkcd(xkcdId)
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, "<b>"+xkcd.Title+"</b>\n───\n<i>"+xkcd.Alt+"</i>\n"+xkcd.Img)
		msg.ParseMode = "HTML"
	}

	bot.Send(msg)
}

func (responder XkcdHandler) Info() *CommandInfo {
	return &xkcdHandlerInfo
}

type XkcdPost struct {
	Title string
	Alt   string
	Img   string
	Id    int
}

var xkcdLatestId int

func GetXkcd(id int) (error, XkcdPost) {
	xkcdStr := "http://xkcd.com/"

	if id != 0 {
		if id == -1 {
			if xkcdLatestId == 0 {
				err, post := GetXkcd(0)
				if err != nil {
					return err, XkcdPost{}
				}
				xkcdLatestId = post.Id
			}
			rand.Seed(time.Now().UTC().UnixNano())
			id = rand.Intn(xkcdLatestId-1) + 1
		}
		xkcdStr += strconv.Itoa(id) + "/"
	}

	xkcdStr += "/info.0.json"

	resp, err := http.Get(xkcdStr)
	if err != nil {
		return err, XkcdPost{}
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Invalid Status Code: " + resp.Status), XkcdPost{}
	}

	var jsonRes map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonRes)
	if err != nil {
		return err, XkcdPost{}
	}

	var xkcdId = int(jsonRes["num"].(float64))
	if xkcdId > xkcdLatestId {
		xkcdLatestId = xkcdId
	}

	return nil, XkcdPost{
		Title: jsonRes["title"].(string),
		Alt:   jsonRes["alt"].(string),
		Img:   jsonRes["img"].(string),
		Id:    xkcdId,
	}

}
