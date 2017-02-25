package main

// GENERATED FILE DO NOT EDIT
// go run tools/gen_cmd_list/main.go -out command_list.go

import (
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"regexp"
)

var CommandHandlers []commands.Command

func LoadCommands() {
	addCommand(&commands.MagicBallHandler{})
	addCommand(&commands.BitcoinHandler{})
	addCommand(&commands.MsgDumpHandler{})
	addCommand(&commands.RedditHandler{})
	addCommand(&commands.RektHandler{})
	addCommand(&commands.ShameHandler{})
	addCommand(&commands.ClearHandler{})
	addCommand(&commands.SettingsHandler{})
	addCommand(&commands.BashHandler{})
	addCommand(&commands.BatmanHandler{})
	addCommand(&commands.BenchHandler{})
	addCommand(&commands.BotFatherHandler{})
	addCommand(&commands.EchoHandler{})
	addCommand(&commands.MemeHandler{})
	addCommand(&commands.StartHandler{})
	addCommand(&commands.UrbanHandler{})
	addCommand(&commands.SquareHandler{})
	addCommand(&commands.MemeListHandler{})
	addCommand(&commands.FortuneHandler{})
	addCommand(&commands.LennyHandler{})
	addCommand(&commands.HelpHandler{})
	addCommand(&commands.InfoHandler{})
	addCommand(&commands.LinesHandler{})
	addCommand(&commands.PingHandler{})
	addCommand(&commands.XkcdHandler{})
	addCommand(&commands.CatHandler{})
	addCommand(&commands.LmgtfyHandler{})
	addCommand(&commands.RandomHandler{})
	addCommand(&commands.VisionHandler{})

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
