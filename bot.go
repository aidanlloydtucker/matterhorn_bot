package main

import (
	"log"

	"strings"

	"strconv"

	"regexp"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/telegram-bot-api.v4"
)

var mainBot *tgbotapi.BotAPI

func startBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	mainBot = bot

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.String(), update.Message.Text)

		go func() {
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
		}()

		log.Println(REGEX_FOR_ALT_COMMAND.MatchString(update.Message.Text))

		if update.Message.Text != "" && (update.Message.IsCommand() || REGEX_FOR_ALT_COMMAND.MatchString(update.Message.Text)) {
			regRes := REGEX_FOR_ALT_COMMAND.FindAllStringSubmatch(update.Message.Text, -1)
			log.Println(regRes)
			if len(regRes) >= 1 && len(regRes[0]) >= 2 {
				update.Message.Text = regRes[0][0]
			}

			for _, cmd := range CommandHandlers {
				if cmd.Info().Command == update.Message.Command() {
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

var REGEX_FOR_ALT_COMMAND *regexp.Regexp = regexp.MustCompilePOSIX(`/<.+>\s\.(.+)/`)

const REDIS_KEY_PREFIX string = "tg-chat-key/"

func formatRedisKey(key int64) string {
	return REDIS_KEY_PREFIX + strconv.FormatInt(key, 10)
}
