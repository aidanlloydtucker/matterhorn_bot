package main

// GENERATED FILE DO NOT EDIT
// go run tools/gen_cmd_list/main.go -out tools/mock_bot/command_list.go

import (
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"regexp"
)

var CommandHandlers []commands.Command

func LoadCommands() {
	addCommand(commands.LinesHandler{})
	addCommand(commands.LmgtfyHandler{})
	addCommand(commands.RandomHandler{})
	addCommand(commands.BotFatherHandler{})
	addCommand(commands.CatHandler{})
	addCommand(commands.EchoHandler{})
	addCommand(commands.LennyHandler{})
	addCommand(commands.RedditHandler{})
	addCommand(commands.RektHandler{})
	addCommand(commands.PingHandler{})
	addCommand(commands.SquareHandler{})
	addCommand(commands.BenchHandler{})
	addCommand(commands.UrbanHandler{})
	addCommand(commands.MagicBallHandler{})
	addCommand(commands.BatmanHandler{})
	addCommand(commands.StartHandler{})
	addCommand(commands.VisionHandler{})
	addCommand(commands.XkcdHandler{})
	addCommand(commands.BitcoinHandler{})
	addCommand(commands.ClearHandler{})
	addCommand(commands.MemeListHandler{})
	addCommand(commands.MsgDumpHandler{})
	addCommand(commands.SettingsHandler{})
	addCommand(commands.ShameHandler{})
	addCommand(commands.BashHandler{})
	addCommand(commands.FortuneHandler{})
	addCommand(commands.HelpHandler{})
	addCommand(commands.InfoHandler{})
	addCommand(commands.MemeHandler{})

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
