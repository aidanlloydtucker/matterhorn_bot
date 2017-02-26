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
	chatpkg "github.com/billybobjoeaglt/matterhorn_bot/chat"
	"github.com/billybobjoeaglt/matterhorn_bot/commands"
	"github.com/codegangsta/cli"
	"google.golang.org/api/option"
)

// Build Vars
var (
	Version   string
	BuildTime string
)

var DatastoreInst *chatpkg.Datastore

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
			Usage: "The IP or domain for the settings webpage and webhook. For webhook, you need this to be a domain ",
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
	log.Println("Running app")

	/* Settings Command Setup */
	HttpPort = c.String("http_port")
	var baseURL string
	if c.IsSet("ip") {
		baseURL = c.String("ip") + ":" + HttpPort
	} else {
		IP, err := checkIP()
		if err != nil {
			baseURL = "localhost:" + HttpPort
		} else {
			baseURL = IP + ":" + HttpPort
		}
	}

	/* Datastore Setup */
	dsCtx := context.Background()

	if !c.IsSet("project_id") {
		log.Fatalln("Missing Project ID")
	}

	// Creates a client.
	dsClient, err := datastore.NewClient(dsCtx, c.String("project_id"), option.WithServiceAccountFile(c.String("service_account_file")))
	if err != nil {
		log.Fatalf("Failed to create datastore client: %v", err)
	}

	DatastoreInst = chatpkg.NewDatastore(dsClient, dsCtx)

	log.Println("Connected to datastore")

	/* Commands */
	LoadCommands()

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
				"url": baseURL + "/chat/",
			})
		case *commands.QuotesHandler:
			cmd.Setup(map[string]interface{}{
				"datastore": DatastoreInst,
			})
		case *commands.QuoteHandler:
			cmd.Setup(map[string]interface{}{
				"datastore": DatastoreInst,
			})
		case *commands.QuoteslinkHandler:
			cmd.Setup(map[string]interface{}{
				"datastore": DatastoreInst,
				"url":       baseURL + "/quotes/",
			})
		default:
			cmd.Setup(map[string]interface{}{})
		}
	}

	log.Println("Loaded all commands")

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
