package main

import (
	"log"

	"strings"

	"strconv"

	"net/http"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/telegram-bot-api.v4"
)

var mainBot *tgbotapi.BotAPI
var runningWebhook bool

type WebhookConfig struct {
	IP       string
	Port     string
	KeyPath  string
	CertPath string
}

func startBot(token string, webhookConf *WebhookConfig) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalln(err)
	}

	mainBot = bot

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	var updates <-chan tgbotapi.Update
	var webhookErr error

	if webhookConf != nil {
		_, webhookErr = bot.SetWebhook(tgbotapi.NewWebhookWithCert(webhookConf.IP+":"+webhookConf.Port+"/"+bot.Token, webhookConf.CertPath))
		if webhookErr != nil {
			log.Println("Webhook Error:", webhookErr, "Switching to poll")
		} else {
			runningWebhook = true
			updates = bot.ListenForWebhook("/" + bot.Token)
			go func() {
				err := http.ListenAndServeTLS("0.0.0.0:"+webhookConf.Port, webhookConf.CertPath, webhookConf.KeyPath, nil)
				if err != nil {
					log.Println("Webhook error:", err)
				}
			}()
			log.Println("Running on Webhook")
		}
	}

	if webhookErr != nil || webhookConf == nil {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err = bot.GetUpdatesChan(u)
		if err != nil {
			log.Fatalln("Error found on getting poll updates:", err, "HALTING")
		}
		log.Println("Running on Poll")
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.String(), update.Message.Text)

		go onMessageRedisRoutine(bot, update)

		if update.Message.Text != "" && update.Message.IsCommand() {
			for _, cmd := range CommandHandlers {
				if cmd.Info().Command == update.Message.Command() {
					if update.Message.ReplyToMessage != nil {
						valid, args := cmd.HandleReply(update.Message)
						if valid {
							update.Message.Text = "/" + cmd.Info().Command + " " + args
						}
					}

					var args []string
					if cmd.Info().Args != "" {
						if cmd.Info().ArgsRegex.MatchString(update.Message.CommandArguments()) {
							matchArr := cmd.Info().ArgsRegex.FindAllStringSubmatch(update.Message.CommandArguments(), -1)
							if len(matchArr) > 0 && len(matchArr[0]) > 1 {
								args = matchArr[0][1:]
							}
						} else {
							continue
						}
					}
					go cmd.HandleCommand(bot, update.Message, args)
				}
			}
		}
	}
}

const REDIS_KEY_PREFIX string = "tg-chat-key/"

func formatRedisKey(key int64) string {
	return REDIS_KEY_PREFIX + strconv.FormatInt(key, 10)
}

func onMessageRedisRoutine(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	exists, err := redis.Bool(redisConn.Do("EXISTS", formatRedisKey(update.Message.Chat.ID)))
	if err != nil {
		return
	}

	if !exists {
		newChat := NewRedisChatInfo()
		if update.Message.Chat.Title == "" {
			if update.Message.Chat.UserName != "" {
				newChat.Name = update.Message.Chat.UserName
			} else {
				newChat.Name = strings.TrimSpace(update.Message.Chat.FirstName + " " + update.Message.Chat.LastName)
			}
		} else {
			newChat.Name = update.Message.Chat.Title
		}
		newChat.Type = update.Message.Chat.Type
		redisConn.Do("HMSET", redis.Args{}.Add(formatRedisKey(update.Message.Chat.ID)).AddFlat(newChat)...)

	} else {
		if update.Message.Text != "" {
			v, err := redis.Values(redisConn.Do("HGETALL", formatRedisKey(update.Message.Chat.ID)))
			if err != nil {
				return
			}

			err, chat := FromRedisChatInfo(v)
			if err != nil {
				return
			}

			for _, word := range chat.KeyWords {
				if strings.Contains(update.Message.Text, word.Key) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, word.Message)
					bot.Send(msg)
				}
			}
		} else if update.Message.NewChatTitle != "" {
			redisConn.Do("HSET", redis.Args{}.Add(formatRedisKey(update.Message.Chat.ID)).Add("name").Add(update.Message.NewChatTitle)...)
		}
	}
}
