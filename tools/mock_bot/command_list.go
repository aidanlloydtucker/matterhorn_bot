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
	addCommand(&commands.BenchHandler{})
	addCommand(&commands.CatHandler{})
	addCommand(&commands.InfoHandler{})
	addCommand(&commands.QuoteslinkHandler{})
	addCommand(&commands.RandomHandler{})
	addCommand(&commands.BashHandler{})
	addCommand(&commands.BotFatherHandler{})
	addCommand(&commands.ClearHandler{})
	addCommand(&commands.FortuneHandler{})
	addCommand(&commands.BitcoinHandler{})
	addCommand(&commands.LennyHandler{})
	addCommand(&commands.ShameHandler{})
	addCommand(&commands.XkcdHandler{})
	addCommand(&commands.MagicBallHandler{})
	addCommand(&commands.QuoteHandler{})
	addCommand(&commands.RedditHandler{})
	addCommand(&commands.StartHandler{})
	addCommand(&commands.MemeHandler{})
	addCommand(&commands.MemeListHandler{})
	addCommand(&commands.MsgDumpHandler{})
	addCommand(&commands.RektHandler{})
	addCommand(&commands.ChangelogHandler{})
	addCommand(&commands.LinesHandler{})
	addCommand(&commands.QuotesHandler{})
	addCommand(&commands.HelpHandler{})
	addCommand(&commands.SquareHandler{})
	addCommand(&commands.VisionHandler{})
	addCommand(&commands.EchoHandler{})
	addCommand(&commands.LmgtfyHandler{})
	addCommand(&commands.PingHandler{})
	addCommand(&commands.SettingsHandler{})
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
