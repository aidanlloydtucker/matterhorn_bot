package main

// GENERATED FILE DO NOT EDIT
// go run tools/gen_cmd_list/main.go -out tools/mock_bot/command_list.go

import (
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"regexp"
)

var CommandHandlers []commands.Command

func LoadCommands() {
	addCommand(&commands.BenchHandler{})
	addCommand(&commands.FortuneHandler{})
	addCommand(&commands.InfoHandler{})
	addCommand(&commands.LinesHandler{})
	addCommand(&commands.VisionHandler{})
	addCommand(&commands.HelpHandler{})
	addCommand(&commands.MsgDumpHandler{})
	addCommand(&commands.ClearHandler{})
	addCommand(&commands.MemeHandler{})
	addCommand(&commands.MemeListHandler{})
	addCommand(&commands.QuotesHandler{})
	addCommand(&commands.RedditHandler{})
	addCommand(&commands.XkcdHandler{})
	addCommand(&commands.BatmanHandler{})
	addCommand(&commands.BashHandler{})
	addCommand(&commands.EchoHandler{})
	addCommand(&commands.LennyHandler{})
	addCommand(&commands.SettingsHandler{})
	addCommand(&commands.CatHandler{})
	addCommand(&commands.PingHandler{})
	addCommand(&commands.ShameHandler{})
	addCommand(&commands.SquareHandler{})
	addCommand(&commands.MagicBallHandler{})
	addCommand(&commands.BitcoinHandler{})
	addCommand(&commands.QuoteHandler{})
	addCommand(&commands.RektHandler{})
	addCommand(&commands.BotFatherHandler{})
	addCommand(&commands.LmgtfyHandler{})
	addCommand(&commands.RandomHandler{})
	addCommand(&commands.StartHandler{})
	addCommand(&commands.UrbanHandler{})

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
