package commands

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"bytes"

	"net/url"

	"strconv"

	"gopkg.in/telegram-bot-api.v4"
)

func newTestBot() (*tgbotapi.BotAPI, chan tgbotapi.MessageConfig) {
	output := make(chan tgbotapi.MessageConfig, 1)
	return &tgbotapi.BotAPI{
		Client: &http.Client{
			Transport: RoundTTest{
				Response: output,
			},
		},
		Token: "foobar",
		Self: tgbotapi.User{
			ID:        1234,
			FirstName: "foo",
			LastName:  "bar",
			UserName:  "foobar",
		},
	}, output
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

type RoundTTest struct {
	Response chan tgbotapi.MessageConfig
}

func (rt RoundTTest) RoundTrip(r *http.Request) (*http.Response, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	vals, err := url.ParseQuery(buf.String())
	if err != nil {
		panic(err)
	}

	rt.Response <- valsToMessageConfig(vals)

	return nil, errors.New("BAD")
}

func valsToMessageConfig(vals url.Values) tgbotapi.MessageConfig {
	msgConf := tgbotapi.MessageConfig{
		Text:      vals.Get("text"),
		ParseMode: vals.Get("parse_mode"),
	}
	dWP, err := strconv.ParseBool(vals.Get("disable_web_page_preview"))
	if err == nil {
		msgConf.DisableWebPagePreview = dWP
	}
	return msgConf
}

func TestBenchHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(BenchHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output

	_, err := strconv.Atoi(out.Text)
	if err != nil {
		t.Error(err)
	}
}

func TestEchoHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	echoStr := "test 123"

	ch := new(EchoHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{echoStr})

	out := <-output

	if out.Text != echoStr {
		t.Fatal("Echo is incorrect")
	}
}
