package main

import (
	"log"

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
