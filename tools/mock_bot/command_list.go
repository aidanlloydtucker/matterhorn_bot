package main

// GENERATED FILE DO NOT EDIT
// go run tools/gen_cmd_list/main.go -out tools/mock_bot/command_list.go

import (
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"regexp"
)

var CommandHandlers []commands.Command

func LoadCommands() {
	addCommand(&commands.InfoHandler{})
	addCommand(&commands.LmgtfyHandler{})
	addCommand(&commands.ShameHandler{})
	addCommand(&commands.UrbanHandler{})
	addCommand(&commands.BitcoinHandler{})
	addCommand(&commands.VisionHandler{})
	addCommand(&commands.StartHandler{})
	addCommand(&commands.BatmanHandler{})
	addCommand(&commands.MemeHandler{})
	addCommand(&commands.PingHandler{})
	addCommand(&commands.SquareHandler{})
	addCommand(&commands.BotFatherHandler{})
	addCommand(&commands.ClearHandler{})
	addCommand(&commands.EchoHandler{})
	addCommand(&commands.LinesHandler{})
	addCommand(&commands.MagicBallHandler{})
	addCommand(&commands.BenchHandler{})
	addCommand(&commands.CatHandler{})
	addCommand(&commands.LennyHandler{})
	addCommand(&commands.RedditHandler{})
	addCommand(&commands.SettingsHandler{})
	addCommand(&commands.XkcdHandler{})
	addCommand(&commands.BashHandler{})
	addCommand(&commands.FortuneHandler{})
	addCommand(&commands.MemeListHandler{})
	addCommand(&commands.MsgDumpHandler{})
	addCommand(&commands.HelpHandler{})
	addCommand(&commands.RandomHandler{})
	addCommand(&commands.RektHandler{})

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
