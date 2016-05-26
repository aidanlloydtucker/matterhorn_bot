package main

import (
	"log"

	"strings"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/telegram-bot-api.v4"
)

func startBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		go func() {
			exists, err := redis.Bool(redisConn.Do("EXISTS", update.Message.Chat.ID))
			if err != nil {
				return
			}

			if !exists {
				newChat := NewChatInfo()
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
				redisConn.Do("HMSET", redis.Args{}.Add(update.Message.Chat.ID).AddFlat(newChat)...)

			} else if update.Message.NewChatTitle != "" {
				redisConn.Do("HSET", redis.Args{}.Add(update.Message.Chat.ID).Add("name").Add(update.Message.NewChatTitle)...)
			}

		}()

		if update.Message.Text != "" && update.Message.IsCommand() {
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
