package custom

import (
	"regexp"

	"github.com/billybobjoeaglt/matterhorn_bot/commands"
)

var CustomCommandList []commands.Command

func addToCustomCommands(cmd commands.Command) {
	if cmd.Info().Args != "" {
		argReg, err := regexp.Compile(cmd.Info().Args)
		if err != nil {
			return
		}
		cmd.Info().ArgsRegex = *argReg
	}

	CustomCommandList = append(CustomCommandList, cmd)
}
