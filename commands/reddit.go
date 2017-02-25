package commands

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"strings"

	"fmt"

	"gopkg.in/telegram-bot-api.v4"
)

type RedditHandler struct {
}

var redditHandlerInfo = CommandInfo{
	Command:     "reddit",
	Args:        `(.[^ ]+) ?(.*)?`,
	Permission:  3,
	Description: "gets a reddit post",
	LongDesc:    "",
	Usage:       "/reddit [subreddit] (sort)",
	Examples: []string{
		"/reddit all",
		"/reddit pics top/all",
		"/reddit funny new",
	},
	ResType: "message",
}

func (h *RedditHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig

	var sort string
	var query string

	if len(args) >= 2 {
		sort = strings.TrimSpace(args[1])
	}

	if strings.Contains(sort, "top") && strings.Contains(sort, "/") {
		split := strings.Split(sort, "/")
		sort = split[0]
		query = "?sort=" + url.QueryEscape(sort) + "&t=" + url.QueryEscape(split[1])
	}

	post, err := GetReddit(bot.Self.String(), args[0], sort, query)
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID,
			fmt.Sprintf("<b>%s</b> - %d points\n──\n%s",
				post.Title, post.Score, post.URL))

		msg.ParseMode = "HTML"
	}
	bot.Send(msg)
}

func (h *RedditHandler) Info() *CommandInfo {
	return &redditHandlerInfo
}

func (h *RedditHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

func (h *RedditHandler) Setup(setupFields map[string]interface{}) {

}

type RedditPost struct {
	Title    string `json:"title"`
	Score    int    `json:"score"`
	URL      string `json:"url"`
	Over18   bool   `json:"over_18"`
	Stickied bool   `json:"stickied"`
}

type RedditResp struct {
	Data struct {
		Children []struct {
			Data RedditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func GetReddit(botName string, subreddit string, sort string, query string) (*RedditPost, error) {
	req, err := http.NewRequest(http.MethodGet,
		"https://www.reddit.com/r/"+url.QueryEscape(subreddit)+
			"/"+url.QueryEscape(sort)+
			".json"+query,
		nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-agent", botName+" 0.1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Invalid Status Code: " + resp.Status)
	}

	redditResp := RedditResp{}
	err = json.NewDecoder(resp.Body).Decode(&redditResp)
	if err != nil {
		return nil, err
	}

	for _, post := range redditResp.Data.Children {
		if !post.Data.Stickied {
			return &post.Data, nil
		}
	}

	return nil, errors.New("Post not Found")

}
