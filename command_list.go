package main

// GENERATED FILE DO NOT EDIT
// go run tools/gen_cmd_list/main.go -out command_list.go

import (
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"regexp"
)

var CommandHandlers []commands.Command

func LoadCommands() {
	addCommand(&commands.BotFatherHandler{})
	addCommand(&commands.QuoteHandler{})
	addCommand(&commands.SettingsHandler{})
	addCommand(&commands.ShameHandler{})
	addCommand(&commands.StartHandler{})
	addCommand(&commands.MagicBallHandler{})
	addCommand(&commands.BashHandler{})
	addCommand(&commands.QuotesHandler{})
	addCommand(&commands.RandomHandler{})
	addCommand(&commands.SquareHandler{})
	addCommand(&commands.VisionHandler{})
	addCommand(&commands.EchoHandler{})
	addCommand(&commands.BenchHandler{})
	addCommand(&commands.LinesHandler{})
	addCommand(&commands.RedditHandler{})
	addCommand(&commands.UrbanHandler{})
	addCommand(&commands.BatmanHandler{})
	addCommand(&commands.ClearHandler{})
	addCommand(&commands.LmgtfyHandler{})
	addCommand(&commands.FortuneHandler{})
	addCommand(&commands.HelpHandler{})
	addCommand(&commands.RektHandler{})
	addCommand(&commands.BitcoinHandler{})
	addCommand(&commands.MemeListHandler{})
	addCommand(&commands.PingHandler{})
	addCommand(&commands.XkcdHandler{})
	addCommand(&commands.InfoHandler{})
	addCommand(&commands.MsgDumpHandler{})
	addCommand(&commands.MemeHandler{})
	addCommand(&commands.CatHandler{})
	addCommand(&commands.LennyHandler{})

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
