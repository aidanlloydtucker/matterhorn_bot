package main

// GENERATED FILE DO NOT EDIT
// go run tools/gen_cmd_list/main.go -out command_list.go

import (
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"regexp"
)

var CommandHandlers []commands.Command

func LoadCommands() {
	addCommand(commands.CatHandler{})
	addCommand(commands.RandomHandler{})
	addCommand(commands.VisionHandler{})
	addCommand(commands.MagicBallHandler{})
	addCommand(commands.HelpHandler{})
	addCommand(commands.LinesHandler{})
	addCommand(commands.BatmanHandler{})
	addCommand(commands.InfoHandler{})
	addCommand(commands.SettingsHandler{})
	addCommand(commands.BitcoinHandler{})
	addCommand(commands.EchoHandler{})
	addCommand(commands.MemeHandler{})
	addCommand(commands.SquareHandler{})
	addCommand(commands.ShameHandler{})
	addCommand(commands.XkcdHandler{})
	addCommand(commands.BotFatherHandler{})
	addCommand(commands.ChatidHandler{})
	addCommand(commands.LennyHandler{})
	addCommand(commands.MemeListHandler{})
	addCommand(commands.PingHandler{})
	addCommand(commands.RektHandler{})
	addCommand(commands.UrbanHandler{})
	addCommand(commands.BashHandler{})
	addCommand(commands.BenchHandler{})
	addCommand(commands.ClearHandler{})
	addCommand(commands.FortuneHandler{})
	addCommand(commands.HotHandler{})
	addCommand(commands.RedditHandler{})
	addCommand(commands.LmgtfyHandler{})
	addCommand(commands.StartHandler{})

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
