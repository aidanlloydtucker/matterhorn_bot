package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"bytes"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/datastore"
	"context"
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"github.com/billybobjoeaglt/matterhorn_bot/commands/custom"
	"github.com/codegangsta/cli"
	"google.golang.org/api/option"
)

// Build Vars
var (
	Version   string
	BuildTime string
)

var datastoreClient *datastore.Client
var datastoreContext context.Context

var HttpPort string

func main() {
	app := cli.NewApp()

	app.Name = "Matterhorn Bot"
	app.Usage = "Telegram bot"

	app.Authors = []cli.Author{
		{
			Name: "Aidan Lloyd-Tucker",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "token, t",
			Usage: "The telegram bot api token",
		},
		cli.StringFlag{
			Name:  "http_port, p",
			Usage: "The http port to open connections for the settings webpage",
			Value: "8080",
		},
		cli.StringFlag{
			Name:  "ip",
			Usage: "The IP for the settings webpage and webhook port",
		},
		cli.StringFlag{
			Name:  "webhook_port",
			Usage: "The telegram bot api webhook port",
			Value: "8443",
		},
		cli.StringFlag{
			Name:  "webhook_cert",
			Usage: "The telegram bot api webhook cert",
			Value: "./ignored/cert.pem",
		},
		cli.StringFlag{
			Name:  "webhook_key",
			Usage: "The telegram bot api webhook key",
			Value: "./ignored/key.key",
		},
		cli.BoolFlag{
			Name:  "enable_webhook, w",
			Usage: "Enables webhook if true",
		},
		cli.BoolFlag{
			Name:  "prod",
			Usage: "Sets bot to production mode",
		},
		cli.StringFlag{
			Name:  "service_account_file",
			Usage: "The filepath of the google service account",
		},
		cli.StringFlag{
			Name:  "project_id",
			Usage: "The project ID for google cloud",
		},
		cli.StringFlag{
			Name:  "set_version",
			Usage: "Set the version of matterhorn bot",
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
	log.Println("Running app")

	HttpPort = c.String("http_port")

	// Add URL for settings command
	var settingsURL string
	if c.IsSet("ip") {
		settingsURL = c.String("ip") + ":" + c.String("http_port") + "/chat/"
	} else {
		IP, err := checkIP()
		if err != nil {
			settingsURL = "localhost:" + c.String("http_port") + "/chat/"
		} else {
			settingsURL = IP + ":" + c.String("http_port") + "/chat/"
		}
	}

	// Commands
	LoadCommands()

	// Load Custom Commands
	custom.LoadCustom()
	for _, cmd := range custom.CustomCommandList {
		CommandHandlers = append(CommandHandlers, cmd)
	}

	cmdMap := make(map[string]*commands.CommandInfo)
	for _, cmd := range CommandHandlers {
		cmdMap[cmd.Info().Command] = cmd.Info()
	}

	for _, cmd := range CommandHandlers {
		switch cmd.(type) {
		case *commands.InfoHandler:
			newVersion := c.String("set_version")
			if newVersion == "" {
				newVersion = c.App.Version
			}
			cmd.Setup(map[string]interface{}{
				"botVersion":   newVersion,
				"botTimestamp": &c.App.Compiled,
			})
		case *commands.HelpHandler:
			cmd.Setup(map[string]interface{}{
				"commandMap": cmdMap,
			})
		case *commands.BotFatherHandler:
			cmd.Setup(map[string]interface{}{
				"commandMap": cmdMap,
			})
		case *commands.VisionHandler:
			cmd.Setup(map[string]interface{}{
				"serviceAccountPath": c.String("service_account_file"),
			})
		case *commands.SettingsHandler:
			cmd.Setup(map[string]interface{}{
				"url": settingsURL,
			})
		default:
			cmd.Setup(map[string]interface{}{})
		}
	}

	log.Println("Loaded all commands")

	/* GOOGLE CLOUD DATASTORE */
	datastoreContext = context.Background()

	if !c.IsSet("project_id") {
		log.Fatalln("Missing Project ID")
	}

	// Creates a client.
	datastoreClient, err = datastore.NewClient(datastoreContext, c.String("project_id"), option.WithServiceAccountFile(c.String("service_account_file")))
	if err != nil {
		log.Fatalf("Failed to create datastore client: %v", err)
	}

	log.Println("Connected to datastore")

	// Start bot

	var webhookConf *WebhookConfig = nil

	if c.IsSet("ip") && c.Bool("enable_webhook") {
		webhookConf = &WebhookConfig{
			IP:       c.String("ip"),
			CertPath: c.String("webhook_cert"),
			KeyPath:  c.String("webhook_key"),
			Port:     c.String("webhook_port"),
		}
	}

	log.Println("Starting bot and website")

	go startBot(c.String("token"), webhookConf)

	// Start Website
	go startWebsite(c.Bool("prod"))

	// Load reminders
	go initTimers()

	// Safe Exit

	var Done = make(chan bool, 1)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs

		if runningWebhook {
			mainBot.RemoveWebhook()
		}

		Done <- true
	}()
	<-Done

	log.Println("Safe Exit")
	return nil
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
