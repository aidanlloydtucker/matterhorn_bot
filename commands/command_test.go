package commands

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

func newLog() *bytes.Buffer {
	var b bytes.Buffer
	//log.SetOutput(&b)

	return &b
}

func newTestBot() *tgbotapi.BotAPI {
	return &tgbotapi.BotAPI{
		Client: &http.Client{},
		Token:  "foobar",
		Self: tgbotapi.User{
			ID:        1234,
			FirstName: "foo",
			LastName:  "bar",
			UserName:  "foobar",
		},
		Debug: true,
	}
}

func newMessageTmpl() *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: 1111,
		From: &tgbotapi.User{
			ID:        2222,
			FirstName: "test",
			LastName:  "user",
			UserName:  "testuser",
		},
		Date: int(time.Now().Unix()),
		Chat: &tgbotapi.Chat{
			ID:    5678,
			Type:  "group",
			Title: "test",
		},
	}
}

func newCommand(command string) *tgbotapi.Message {
	msg := newMessageTmpl()
	msg.Text = "/" + command
	return msg
}

func TestBenchHandler_HandleCommand(t *testing.T) {
	out := newLog()
	bot := newTestBot()

	bh := new(BenchHandler)
	bh.HandleCommand(bot, newCommand(bh.Info().Command), []string{})
	t.Log("OUT", out.Bytes())
	//endpoint, params, message
}
