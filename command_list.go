package main

// GENERATED FILE DO NOT EDIT
// go run tools/gen_cmd_list/main.go -out command_list.go

import (
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"regexp"
)

var CommandHandlers []commands.Command

func LoadCommands() {
	addCommand(&commands.LinesHandler{})
	addCommand(&commands.MsgDumpHandler{})
	addCommand(&commands.QuotesHandler{})
	addCommand(&commands.QuoteslinkHandler{})
	addCommand(&commands.RedditHandler{})
	addCommand(&commands.ShameHandler{})
	addCommand(&commands.StartHandler{})
	addCommand(&commands.BitcoinHandler{})
	addCommand(&commands.LmgtfyHandler{})
	addCommand(&commands.UrbanHandler{})
	addCommand(&commands.QuoteHandler{})
	addCommand(&commands.FortuneHandler{})
	addCommand(&commands.LennyHandler{})
	addCommand(&commands.CatHandler{})
	addCommand(&commands.InfoHandler{})
	addCommand(&commands.RektHandler{})
	addCommand(&commands.MagicBallHandler{})
	addCommand(&commands.BashHandler{})
	addCommand(&commands.MemeHandler{})
	addCommand(&commands.RandomHandler{})
	addCommand(&commands.ClearHandler{})
	addCommand(&commands.SquareHandler{})
	addCommand(&commands.XkcdHandler{})
	addCommand(&commands.BotFatherHandler{})
	addCommand(&commands.ChangelogHandler{})
	addCommand(&commands.MemeListHandler{})
	addCommand(&commands.BatmanHandler{})
	addCommand(&commands.HelpHandler{})
	addCommand(&commands.PingHandler{})
	addCommand(&commands.SettingsHandler{})
	addCommand(&commands.VisionHandler{})
	addCommand(&commands.BenchHandler{})
	addCommand(&commands.EchoHandler{})

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
