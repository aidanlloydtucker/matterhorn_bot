package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"github.com/billybobjoeaglt/matterhorn_bot/commands/custom"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	// Commands
	LoadCommands()

	// Load Custom Commands
	custom.LoadCustom()
	for _, cmd := range custom.CustomCommandList {
		CommandHandlers = append(CommandHandlers, cmd)
	}

	cmdMap := make(map[string]*commands.CommandInfo)
	for _, cmd := range CommandHandlers {
		cmdMap[cmd.Info().Command] = cmd.Info()
	}

	// Help Command Setup
	commands.CommandMap = cmdMap

	// ACTUAL

	reader := bufio.NewReader(os.Stdin)
	bot, botOut := newTestBot()
	for {
		fmt.Print("Enter Command: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		metaCommandArgsMap := map[string]string{}

		// [PHOTO="path/to/photo" REPLY="hello world"]
		if strings.HasPrefix(text, "[") {
			var inner string
			var commandText string

			_, err = fmt.Sscanf(text, "[%s] %s", &inner, &commandText)
			if err != nil {
				panic(err)
			}

			text = commandText



			argsSlice := strings.Split(inner, `"`)
			for i, arg := range argsSlice {
				arg = strings.TrimSpace(arg)
				if i%2 == 0 && len(argsSlice) > i+1 {
					metaCommandArgsMap[arg] = strings.TrimSpace(argsSlice[i+1])
				}
			}
		}

		msg := newInputMessage(text)

		commandSent := false

		for _, cmd := range CommandHandlers {
			if cmd.Info().Command == msg.Command() {
				if cmd.Info().ResType != "message" {
					fmt.Println("Cannot use a command that doesnt return a message text")
					break
				}

				var args []string
				if cmd.Info().Args != "" {
					if cmd.Info().ArgsRegex.MatchString(msg.CommandArguments()) {
						matchArr := cmd.Info().ArgsRegex.FindAllStringSubmatch(msg.CommandArguments(), -1)
						if len(matchArr) > 0 && len(matchArr[0]) > 1 {
							args = matchArr[0][1:]
						}
					} else {
						continue
					}
				}

				cmd.HandleCommand(bot, msg, args)
				commandSent = true
				break
			}
		}
		if commandSent {
			out := <-botOut
			fmt.Println(strings.TrimSpace(out.Text) + "\n")
		} else {
			fmt.Println("Unknown Command")
		}
	}
}
func newInputMessage(text string) *tgbotapi.Message {
	text = strings.TrimSpace(text)
	return &tgbotapi.Message{
		MessageID: 1,
		From: &tgbotapi.User{
			ID:        2,
			FirstName: "Charlie",
			LastName:  "Brown",
			UserName:  "charliebrown",
		},
		Date: int(time.Now().Unix()),
		Chat: &tgbotapi.Chat{
			ID:    3,
			Type:  "group",
			Title: "TV",
		},
		Text: text,
	}
}

func AddCommand(cmd commands.Command) {
	if cmd.Info().Args != "" {
		argReg, err := regexp.Compile(cmd.Info().Args)
		if err != nil {
			return
		}
		cmd.Info().ArgsRegex = *argReg
	}

	CommandHandlers = append(CommandHandlers, cmd)
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
