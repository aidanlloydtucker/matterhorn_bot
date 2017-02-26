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
	addCommand(&commands.BenchHandler{})
	addCommand(&commands.InfoHandler{})
	addCommand(&commands.QuoteslinkHandler{})
	addCommand(&commands.RektHandler{})
	addCommand(&commands.LinesHandler{})
	addCommand(&commands.MemeHandler{})
	addCommand(&commands.BashHandler{})
	addCommand(&commands.BatmanHandler{})
	addCommand(&commands.SettingsHandler{})
	addCommand(&commands.XkcdHandler{})
	addCommand(&commands.EchoHandler{})
	addCommand(&commands.LmgtfyHandler{})
	addCommand(&commands.ShameHandler{})
	addCommand(&commands.StartHandler{})
	addCommand(&commands.BotFatherHandler{})
	addCommand(&commands.HelpHandler{})
	addCommand(&commands.BitcoinHandler{})
	addCommand(&commands.FortuneHandler{})
	addCommand(&commands.MsgDumpHandler{})
	addCommand(&commands.UrbanHandler{})
	addCommand(&commands.VisionHandler{})
	addCommand(&commands.RandomHandler{})
	addCommand(&commands.QuoteHandler{})
	addCommand(&commands.QuotesHandler{})
	addCommand(&commands.CatHandler{})
	addCommand(&commands.ClearHandler{})
	addCommand(&commands.LennyHandler{})
	addCommand(&commands.MemeListHandler{})
	addCommand(&commands.PingHandler{})
	addCommand(&commands.RedditHandler{})
	addCommand(&commands.SquareHandler{})

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
