package main

// GENERATED FILE DO NOT EDIT
// go run tools/gen_cmd_list/main.go -out command_list.go

import (
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"regexp"
)

var CommandHandlers []commands.Command

func LoadCommands() {
	addCommand(commands.EchoHandler{})
	addCommand(commands.LinesHandler{})
	addCommand(commands.BenchHandler{})
	addCommand(commands.CatHandler{})
	addCommand(commands.FortuneHandler{})
	addCommand(commands.LmgtfyHandler{})
	addCommand(commands.SettingsHandler{})
	addCommand(commands.BatmanHandler{})
	addCommand(commands.BotFatherHandler{})
	addCommand(commands.ClearHandler{})
	addCommand(commands.InfoHandler{})
	addCommand(commands.RandomHandler{})
	addCommand(commands.SquareHandler{})
	addCommand(commands.UrbanHandler{})
	addCommand(commands.BitcoinHandler{})
	addCommand(commands.BashHandler{})
	addCommand(commands.LennyHandler{})
	addCommand(commands.PingHandler{})
	addCommand(commands.ShameHandler{})
	addCommand(commands.VisionHandler{})
	addCommand(commands.MagicBallHandler{})
	addCommand(commands.MemeListHandler{})
	addCommand(commands.StartHandler{})
	addCommand(commands.HotHandler{})
	addCommand(commands.MemeHandler{})
	addCommand(commands.XkcdHandler{})
	addCommand(commands.HelpHandler{})
	addCommand(commands.RedditHandler{})
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
