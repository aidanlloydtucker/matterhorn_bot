package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/billybobjoeaglt/sansa_bot/commands"
	"github.com/codegangsta/cli"
)

var (
	BotToken  string
	Version   string
	BuildTime string
)

var CommandHandlers []commands.Command

func main() {
	app := cli.NewApp()

	app.Name = "AutoMod Bot"
	app.Usage = "Telegram bot"

	app.Authors = []cli.Author{
		cli.Author{
			Name: "Aidan Lloyd-Tucker",
		},
	}

	app.Version = Version

	num, err := strconv.ParseInt(BuildTime, 10, 64)
	if err == nil {
		app.Compiled = time.Unix(num, 0)
	}

	app.Action = runApp
	app.Run(os.Args)
}

func runApp(c *cli.Context) {

	// Commands

	AddCommand(commands.BatmanHandler{})
	AddCommand(commands.BenchHandler{})
	AddCommand(commands.BitcoinHandler{})
	AddCommand(commands.CatHandler{})
	AddCommand(commands.UrbanHandler{})
	AddCommand(commands.ClearHandler{})
	AddCommand(commands.EchoHandler{})
	AddCommand(commands.HelpHandler{})
	AddCommand(commands.FortuneHandler{})
	AddCommand(commands.LennyHandler{})

	// Help Command Setup
	commands.CommandList = &CommandHandlers

	// Start bot

	startBot(BotToken)

	// Safe Exit

	var Done = make(chan bool, 1)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		Done <- true
	}()
	<-Done

	log.Println("Safe Exit")
}

func AddCommand(cmd commands.Command) {
	if cmd.Info().Args != "" {
		argReg, err := regexp.Compile(cmd.Info().Args)
		if err != nil {
			return
		}
		cmd.Info().ArgsRegex = *argReg
	}

	CommandHandlers = append(CommandHandlers, cmd)
}
