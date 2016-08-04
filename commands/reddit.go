package commands

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"strings"

	"strconv"

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
		"/reddit funny",
		"/reddit pics top/all",
		"/reddit golang new",
	},
	ResType: "message",
}

func (h RedditHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
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

	err, post := GetReddit(bot.Self.String(), args[0], sort, query)
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, "<b>"+post.Title+"</b> - "+post.Score+" points\n───\n"+post.URL+"")
		msg.ParseMode = "HTML"
	}
	bot.Send(msg)
}

func (h RedditHandler) Info() *CommandInfo {
	return &redditHandlerInfo
}

func (h RedditHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

type RedditPost struct {
	Title string
	Score string
	URL   string
}

func GetReddit(botName string, subreddit string, sort string, query string) (error, RedditPost) {
	req, err := http.NewRequest(http.MethodGet, "https://www.reddit.com/r/"+url.QueryEscape(subreddit)+"/"+url.QueryEscape(sort)+".json"+query, nil)
	if err != nil {
		return err, RedditPost{}
	}
	req.Header.Add("User-agent", botName+" 0.1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err, RedditPost{}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Invalid Status Code: " + resp.Status), RedditPost{}
	}

	var jsonRes map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonRes)
	if err != nil {
		return err, RedditPost{}
	}

	if len(jsonRes["data"].(map[string]interface{})) <= 0 {
		return errors.New("Missing JSON Data"), RedditPost{}
	}

	postList := jsonRes["data"].(map[string]interface{})["children"].([]interface{})
	if len(postList) <= 0 {
		return errors.New("Missing Post"), RedditPost{}
	}
	for _, post := range postList {
		postData := post.(map[string]interface{})["data"].(map[string]interface{})
		if !postData["stickied"].(bool) {
			return nil, RedditPost{
				Title: postData["title"].(string),
				Score: strconv.FormatFloat(postData["score"].(float64), 'f', 0, 64),
				URL:   postData["url"].(string),
			}
		}
	}

	return errors.New("Post not Found"), RedditPost{}

}
