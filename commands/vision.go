package commands

import (
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"context"
	"cloud.google.com/go/vision"
	"log"
	"fmt"
	"strings"
	"google.golang.org/api/option"
)

type VisionHandler struct {
}

var visionHandlerInfo = CommandInfo{
	Command:     "vision",
	Args:        "",
	Permission:  3,
	Description: "gets image labels from google vision ai",
	LongDesc:    "",
	Usage:       "/vision",
	Examples: []string{
		"/vision",
	},
	ResType: "message",
}

var visionClient *vision.Client
var visionContext context.Context

func LoadVision(serviceAcctPath string) {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewClient(ctx, option.WithServiceAccountFile(serviceAcctPath))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	} else {
		visionClient = client
		visionContext = ctx
	}
}

func (h VisionHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var msg tgbotapi.MessageConfig
	var errMsg tgbotapi.MessageConfig
	var err error

	defer func(bot *tgbotapi.BotAPI) {
		if errMsg.Text != "" {
			bot.Send(errMsg)
		}
	}(bot)

	replyToMsg := message.MessageID

	if message.Photo == nil {
		if message.ReplyToMessage != nil && message.ReplyToMessage.Photo != nil {
			message.Photo = message.ReplyToMessage.Photo
			replyToMsg = message.ReplyToMessage.MessageID
		} else {
			return
		}
	}

	fileID := (*(message.Photo))[len(*(message.Photo))-1].FileID

	vis, err := getVision(bot, fileID)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return
	}

	labels := make([]string, len(vis.Labels))
	for i, label := range vis.Labels {
		labels[i] = "• " + label + "\n"
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "<b>Labels</b>\n" + strings.Join(labels, ""))
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = replyToMsg

	_, err = bot.Send(msg)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
	}


}

func (h VisionHandler) Info() *CommandInfo {
	return &visionHandlerInfo
}

func (h VisionHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return true, ""
}

type Vision struct {
	Labels []string
}

func getVision(bot *tgbotapi.BotAPI, fileID string) (*Vision, error) {
	fileurl, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		return nil, err
	}

	photoResp, err := http.Get(fileurl)
	if err != nil {
		return nil, err
	}

	image, err := vision.NewImageFromReader(photoResp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to create image: %v", err)
	}

	labels, err := visionClient.DetectLabels(visionContext, image, 10)
	if err != nil {
		return nil, fmt.Errorf("Failed to detect labels: %v", err)
	}

	labs := []string{}

	for _, label := range labels {
		labs = append(labs, label.Description)
	}
	return &Vision{labs}, nil
}
