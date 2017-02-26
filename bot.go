package main

import (
	"log"

	"strings"

	"net/http"

	"gopkg.in/telegram-bot-api.v4"

	chatpkg "github.com/billybobjoeaglt/matterhorn_bot/chat"
	mbCommands "github.com/billybobjoeaglt/matterhorn_bot/commands"
	"math/rand"
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
	log.Println("Starting Bot")

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
		var cert interface{} = nil
		if webhookConf.CertPath != "" {
			cert = webhookConf.CertPath
		}
		_, webhookErr = bot.SetWebhook(tgbotapi.NewWebhookWithCert(webhookConf.IP+":"+webhookConf.Port+"/"+bot.Token, cert))
		if webhookErr != nil {
			log.Println("Webhook Error:", webhookErr, "Switching to poll")
		} else {
			runningWebhook = true
			updates = bot.ListenForWebhook("/" + bot.Token)
			go http.ListenAndServeTLS(webhookConf.IP+":"+webhookConf.Port, webhookConf.CertPath, webhookConf.KeyPath, nil)
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

		if update.Message.Caption != "" && update.Message.Text == "" {
			update.Message.Text = update.Message.Caption
		}

		log.Printf("[%s] %s", update.Message.From.String(), update.Message.Text)

		go onMessageRoutine(bot, update)

		// FOR INLINE 8BALL COMMANDS. FIGURE OUT A BETTER, MORE MODULAR WAY TO DO THIS LATER
		//TODO: Make this less jankey and make it support Setup()
		if update.Message.Text != "" && strings.Contains(update.Message.Text, "#8ball") {
			go (&mbCommands.MagicBallHandler{}).HandleCommand(bot, update.Message, []string{})
		}
		// ENDING THE 8BALL COMMAND SECTION

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

func onMessageRoutine(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chat, exists, err := DatastoreInst.GetChat(update.Message.Chat.ID)
	if err != nil && exists {
		log.Println("Error getting chat from datastore:", err)
		return
	} else if !exists {
		newChat := chatpkg.Chat{
			Settings: chatpkg.ChatSettings{
				QuotesDoc: int(rand.Int31()),
			},
		}
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

		err = DatastoreInst.InsertChat(newChat, update.Message.Chat.ID)
		if err != nil {
			log.Println("Error inserting chat:", err)
		} else {
			log.Printf("Successfully inserted chat %v (%v)\n", chat.Name, update.Message.Chat.ID)
		}

	} else {
		if chat.Settings.QuotesDoc == 0 {
			_, err = DatastoreInst.UpdateChat(func(oldChat chatpkg.Chat) chatpkg.Chat {
				newChat := oldChat
				newChat.Settings.QuotesDoc = int(rand.Int31())
				return newChat
			}, update.Message.Chat.ID)
			if err != nil {
				log.Println("Error updating quotesdoc chat:", err)
			}
		}

		if update.Message.Text != "" {
			for _, word := range chat.Settings.KeyWords {
				if strings.Contains(strings.ToLower(update.Message.Text), strings.ToLower(word.Key)) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, word.Message)
					bot.Send(msg)
				}
			}
		} else if update.Message.NewChatTitle != "" {
			chat.Name = update.Message.NewChatTitle
			_, err = DatastoreInst.UpdateChat(func(oldChat chatpkg.Chat) chatpkg.Chat {
				newChat := oldChat
				newChat.Name = update.Message.NewChatTitle
				return newChat
			}, update.Message.Chat.ID)
			if err != nil {
				log.Println("Error updating chat:", err)
			}
		}
	}
}
