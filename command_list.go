package main

// GENERATED FILE DO NOT EDIT
// go run tools/gen_cmd_list/main.go -out command_list.go

import (
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"regexp"
)

var CommandHandlers []commands.Command

func LoadCommands() {
	addCommand(commands.XkcdHandler{})
	addCommand(commands.RedditHandler{})
	addCommand(commands.SquareHandler{})
	addCommand(commands.SettingsHandler{})
	addCommand(commands.BatmanHandler{})
	addCommand(commands.PingHandler{})
	addCommand(commands.LinesHandler{})
	addCommand(commands.MsgDumpHandler{})
	addCommand(commands.BashHandler{})
	addCommand(commands.BitcoinHandler{})
	addCommand(commands.RandomHandler{})
	addCommand(commands.MemeListHandler{})
	addCommand(commands.BotFatherHandler{})
	addCommand(commands.MemeHandler{})
	addCommand(commands.LennyHandler{})
	addCommand(commands.LmgtfyHandler{})
	addCommand(commands.VisionHandler{})
	addCommand(commands.BenchHandler{})
	addCommand(commands.InfoHandler{})
	addCommand(commands.StartHandler{})
	addCommand(commands.MagicBallHandler{})
	addCommand(commands.EchoHandler{})
	addCommand(commands.FortuneHandler{})
	addCommand(commands.HelpHandler{})
	addCommand(commands.RektHandler{})
	addCommand(commands.ShameHandler{})
	addCommand(commands.UrbanHandler{})
	addCommand(commands.CatHandler{})
	addCommand(commands.ClearHandler{})

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
