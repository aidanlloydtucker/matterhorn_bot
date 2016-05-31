package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/billybobjoeaglt/sansa_bot/commands"
	"github.com/codegangsta/cli"
	"github.com/garyburd/redigo/redis"
)

var (
	BotToken  string
	Version   string
	BuildTime string
	HttpPort  string
	IP        string
)

var redisConn redis.Conn

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

func runApp(c *cli.Context) error {
	var err error

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
	AddCommand(commands.BashHandler{})
	AddCommand(commands.LmgtfyHandler{})
	AddCommand(commands.PingHandler{})
	AddCommand(commands.RedditHandler{})
	AddCommand(commands.LinesHandler{})
	AddCommand(commands.SquareHandler{})
	AddCommand(commands.StartHandler{})
	AddCommand(commands.XkcdHandler{})
	AddCommand(commands.BotFatherHandler{})
	AddCommand(commands.SettingsHandler{})

	// Help Command Setup
	commands.CommandList = &CommandHandlers

	// Connect to redis
	redisConn, err = redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	defer redisConn.Close()

	// Set HttpPort
	if HttpPort == "" {
		HttpPort = "8080"
	}

	// Add URL for settings command
	if IP != "" {
		commands.SettingsURL = IP + ":" + HttpPort + "/chat/"
	} else {
		ip, err := checkIP()
		if err != nil {
			commands.SettingsURL = "localhost:" + HttpPort + "/chat/"
		} else {
			commands.SettingsURL = ip + ":" + HttpPort + "/chat/"
		}
	}

	// Start bot

	go startBot(BotToken)

	// Start Website

	go startWebsite()

	// Load reminders
	go loadTimeReminders()

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
	return nil
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

func checkIP() (string, error) {
	rsp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(buf)), nil
}
