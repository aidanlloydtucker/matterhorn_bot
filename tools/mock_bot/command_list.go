package main

// GENERATED FILE DO NOT EDIT
// go run tools/gen_cmd_list/main.go -out tools/mock_bot/command_list.go

import (
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"regexp"
)

var CommandHandlers []commands.Command

func LoadCommands() {
	addCommand(&commands.BatmanHandler{})
	addCommand(&commands.BotFatherHandler{})
	addCommand(&commands.LinesHandler{})
	addCommand(&commands.RedditHandler{})
	addCommand(&commands.RektHandler{})
	addCommand(&commands.MagicBallHandler{})
	addCommand(&commands.BashHandler{})
	addCommand(&commands.LennyHandler{})
	addCommand(&commands.RandomHandler{})
	addCommand(&commands.VisionHandler{})
	addCommand(&commands.QuoteslinkHandler{})
	addCommand(&commands.XkcdHandler{})
	addCommand(&commands.CatHandler{})
	addCommand(&commands.LmgtfyHandler{})
	addCommand(&commands.MemeListHandler{})
	addCommand(&commands.QuoteHandler{})
	addCommand(&commands.InfoHandler{})
	addCommand(&commands.MsgDumpHandler{})
	addCommand(&commands.QuotesHandler{})
	addCommand(&commands.BenchHandler{})
	addCommand(&commands.MemeHandler{})
	addCommand(&commands.StartHandler{})
	addCommand(&commands.BitcoinHandler{})
	addCommand(&commands.PingHandler{})
	addCommand(&commands.ShameHandler{})
	addCommand(&commands.SquareHandler{})
	addCommand(&commands.ClearHandler{})
	addCommand(&commands.HelpHandler{})
	addCommand(&commands.SettingsHandler{})
	addCommand(&commands.UrbanHandler{})
	addCommand(&commands.EchoHandler{})
	addCommand(&commands.FortuneHandler{})

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
