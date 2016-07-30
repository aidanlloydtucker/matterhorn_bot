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

	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"github.com/billybobjoeaglt/matterhorn_bot/commands/custom"
	"github.com/codegangsta/cli"
	"github.com/garyburd/redigo/redis"
)

var (
	BotToken      string
	Version       string
	BuildTime     string
	HttpPort      string
	IP            string
	WebhookPort   string
	WebhookCert   string
	WebhookKey    string
	EnableWebhook string
)

var redisConn redis.Conn

var CommandHandlers []commands.Command

func main() {
	app := cli.NewApp()

	app.Name = "Matterhorn Bot"
	app.Usage = "Telegram bot"

	app.Authors = []cli.Author{
		{
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

func setDefaults() {
	if HttpPort == "" {
		HttpPort = "8080"
	}

	if WebhookPort == "" {
		WebhookPort = "8443"
	}

	if WebhookCert == "" {
		WebhookCert = "./ignored/cert.pem"
	}

	if WebhookKey == "" {
		WebhookKey = "./ignored/key.key"
	}
}

func runApp(c *cli.Context) error {
	setDefaults()

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
	AddCommand(commands.MemeHandler{})
	AddCommand(commands.MemeListHandler{})
	AddCommand(commands.ShameHandler{})

	// Load Custom Commands
	custom.LoadCustom()
	for _, cmd := range custom.CustomCommandList {
		CommandHandlers = append(CommandHandlers, cmd)
	}

	// Help Command Setup
	commands.CommandList = &CommandHandlers

	// Connect to redis
	redisConn, err = redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	defer redisConn.Close()

	// Add URL for settings command
	if IP != "" {
		commands.SettingsURL = IP + ":" + HttpPort + "/chat/"
	} else {
		IP, err = checkIP()
		if err != nil {
			commands.SettingsURL = "localhost:" + HttpPort + "/chat/"
		} else {
			commands.SettingsURL = IP + ":" + HttpPort + "/chat/"
		}
	}

	// Start bot

	var webhookConf *WebhookConfig = nil

	if IP != "" && EnableWebhook == "YES" {
		webhookConf = &WebhookConfig{
			IP:       IP,
			CertPath: WebhookCert,
			KeyPath:  WebhookKey,
			Port:     WebhookPort,
		}
	}

	go startBot(BotToken, webhookConf)

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
