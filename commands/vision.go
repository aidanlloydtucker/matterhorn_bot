package commands

import (
	"cloud.google.com/go/vision"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"net/http"
	"strings"
)

type VisionHandler struct {
	client    *vision.Client
	clientCtx context.Context
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

func (h *VisionHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
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

	vis, err := getVision(bot, h.client, h.clientCtx, fileID)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return
	}

	labels := make([]string, len(vis.Labels))
	for i, label := range vis.Labels {
		labels[i] = "• " + label + "\n"
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "<b>Labels</b>\n"+strings.Join(labels, ""))
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = replyToMsg

	_, err = bot.Send(msg)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
	}

}

func (h *VisionHandler) Info() *CommandInfo {
	return &visionHandlerInfo
}

func (h *VisionHandler) HandleReply(message *tgbotapi.Message) (bool, string) {
	return true, ""
}

/*
Params:
string serviceAccountPath (default: "") // path to the google service account file which authorizes access to the vision API. if blank (""), the command will not work
*/
func (h *VisionHandler) Setup(setupFields map[string]interface{}) {
	var serviceAcctPath string

	if pathVal, ok := setupFields["serviceAccountPath"]; ok {
		if newPath, ok := pathVal.(string); ok {
			serviceAcctPath = newPath
		}
	}

	if serviceAcctPath == "" {
		log.Println("Failed to create vision client: missing service account path")
		return
	}

	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewClient(ctx, option.WithServiceAccountFile(serviceAcctPath))
	if err != nil {
		log.Printf("Failed to create vision client: %v\n", err)
	} else {
		h.client = client
		h.clientCtx = ctx
	}
}

type Vision struct {
	Labels []string
}

func getVision(bot *tgbotapi.BotAPI, client *vision.Client, clientCtx context.Context, fileID string) (*Vision, error) {
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

	labels, err := client.DetectLabels(clientCtx, image, 10)
	if err != nil {
		return nil, fmt.Errorf("Failed to detect labels: %v", err)
	}

	labs := []string{}

	for _, label := range labels {
		labs = append(labs, label.Description)
	}
	return &Vision{labs}, nil
}
