package commands

import "regexp"

type CommandInfo struct {
	Command     string
	Args        string
	ArgsRegex   regexp.Regexp
	Permission  int
	Description string
	LongDesc    string
	Usage       string
	Examples    []string
	ResType     string
}
