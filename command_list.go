package main

// GENERATED FILE DO NOT EDIT
// go run tools/gen_cmd_list/main.go --out command_list.go

import (
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"regexp"
)

var CommandHandlers []commands.Command

func LoadCommands() {
	addCommand(commands.MagicBallHandler{})
	addCommand(commands.BashHandler{})
	addCommand(commands.CatHandler{})
	addCommand(commands.ShameHandler{})
	addCommand(commands.EchoHandler{})
	addCommand(commands.MemeListHandler{})
	addCommand(commands.SquareHandler{})
	addCommand(commands.HotHandler{})
	addCommand(commands.PingHandler{})
	addCommand(commands.RedditHandler{})
	addCommand(commands.StartHandler{})
	addCommand(commands.InfoHandler{})
	addCommand(commands.SettingsHandler{})
	addCommand(commands.XkcdHandler{})
	addCommand(commands.BatmanHandler{})
	addCommand(commands.LennyHandler{})
	addCommand(commands.LmgtfyHandler{})
	addCommand(commands.MemeHandler{})
	addCommand(commands.RandomHandler{})
	addCommand(commands.UrbanHandler{})
	addCommand(commands.BenchHandler{})
	addCommand(commands.BitcoinHandler{})
	addCommand(commands.ClearHandler{})
	addCommand(commands.HelpHandler{})
	addCommand(commands.BotFatherHandler{})
	addCommand(commands.FortuneHandler{})
	addCommand(commands.LinesHandler{})
	addCommand(commands.RektHandler{})

}

func addCommand(cmd commands.Command) {
	if cmd.Info().Args != "" {
		argReg, err := regexp.Compile(cmd.Info().Args)
		if err != nil {
			return
		}
		cmd.Info().ArgsRegex = *argReg
	}

	CommandHandlers = append(CommandHandlers, cmd)
}
