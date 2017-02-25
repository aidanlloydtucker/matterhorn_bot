package commands

import (
	"errors"
	"fmt"

	"gopkg.in/telegram-bot-api.v4"
)

type HelpHandler struct {
	commandMap map[string]*CommandInfo
}

var helpHandlerInfo = CommandInfo{
	Command:     "help",
	Args:        "",
	Permission:  3,
	Description: "lists commands",
	LongDesc:    "",
	Usage:       "/help (command)",
	Examples: []string{
		"/help",
		"/help ping",
	},
	ResType: "message",
}

func (h *HelpHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msgStr string
	var err error
	if message.CommandArguments() != "" {
		msgStr, err = getCommandInfo(h.commandMap, message.CommandArguments())
	} else {
		msgStr = listCommands(h.commandMap)
	}

	var msg tgbotapi.MessageConfig
	if err != nil {
		msg = NewErrorMessage(message.Chat.ID, err)
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, msgStr)
		msg.ParseMode = "HTML"
	}
	bot.Send(msg)
}

func listCommands(cmdMap map[string]*CommandInfo) string {
	msgStr := "<b>Commands:</b>\n"
	for _, cmd := range cmdMap {
		if cmd.Hidden {
			continue
		}
		msgStr += "• " + cmd.Command + " - " + cmd.Description + "\n"
	}
	return msgStr
}

func getCommandInfo(cmdMap map[string]*CommandInfo, command string) (string, error) {
	cmd, ok := cmdMap[command]
	if !ok {
		return "", errors.New("Unknown command")
	}

	desc := cmd.LongDesc
	if desc == "" {
		desc = cmd.Description
	}

	var examples string
	for _, ex := range cmd.Examples {
		examples += fmt.Sprintf("\t• <code>%s</code>\n", ex)
	}

	msgStr := fmt.Sprintf(`<b>%s</b>
———
<b>Description</b> %s. Will return a %s.
<b>Usage</b> <pre>%s</pre>
<b>Examples</b>
%v`, cmd.Command, desc, cmd.ResType, cmd.Usage, examples)

	return msgStr, nil
}

func (h *HelpHandler) Info() *CommandInfo {
	return &helpHandlerInfo
}

func (h *HelpHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return false, ""
}

/*
Params:
map[string]*CommandInfo commandMap (default: map[string]*CommandInfo{}) // Map of command infos for help to print out
*/
func (h *HelpHandler) Setup(setupFields map[string]interface{}) {
	h.commandMap = map[string]*CommandInfo{}

	if cmdMapVal, ok := setupFields["commandMap"]; ok {
		if cmdMap, ok := cmdMapVal.(map[string]*CommandInfo); ok {
			h.commandMap = cmdMap
		}
	}
}
